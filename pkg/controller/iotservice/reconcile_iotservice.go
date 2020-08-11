package iotservice

import (
	"context"
	"fmt"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	cpv1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1"
	cpv1alpha1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/configmap"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/deployment"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/namespace"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/network_policy"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/secret"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/service"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/stateful_set"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var _ reconcile.Reconciler = &ReconcileIoTService{}
var currentUpdateCount = 0

// ReconcileIoTService reconciles a IoTService object
type ReconcileIoTService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile is the main entry point for runtime of the Operator.
// As soon as something has changed on the specification or the status of a watched resource,
// this method is called to align the as-is with the to-be state.
// The method involves
// - handling the finalizer of dropping namespace upon instance deletion
// - limiting the number of parallel updates to save cluster resources
// - updating resources according to to-be state
// - updating status of IoTService object
func (r *ReconcileIoTService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	logger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	logger.Info("Reconciling IoTService")

	// requests from the primary resource come from the operator namespace,
	// requests from the secondary resources come from the instance namespace.
	// Fetch the IoTService instance from the operator namespace
	operatorNamespace, err := k8sutil.GetOperatorName()
	if err != nil {
		return reconcile.Result{}, err
	}

	instance := &latest.IoTService{}

	instanceV1 := &cpv1.IoTService{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: request.Namespace, Name: request.Name}, instanceV1)

	if err != nil {

			instanceV1alpha1 := &cpv1alpha1.IoTService{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Namespace: request.Namespace, Name: request.Name}, instanceV1alpha1)

			if err != nil {
				if errors.IsNotFound(err) {
					// Request object not found:
					// a) could have been deleted after reconcile request.
					// b) the request could be of a secondary resource while the primary resource has been deleted
					// Return and don't requeue
					logger.Info("IoTService object not found. Remaining resources are garbage collected.")
					return reconcile.Result{}, nil
				}
				// Error reading the object - requeue the request.
				return reconcile.Result{}, err
			} else {
				logger.Info("IoTService object found, version v1alpha1")
				instance.V1alpha1 = instanceV1alpha1
				instance.SetInstanceId(instanceV1alpha1.Name)
			}

	} else {
		logger.Info("IoTService object found, version v1")
		latest.FillDefaultValues(instanceV1)
		instance.V1 = instanceV1
		instance.SetInstanceId(instanceV1.GetName())
	}

	instance.SetNamespace(request.Namespace)

	requeue := false

	var isDeletion bool
	if err := r.handleFinalizer(instance, "com.sap.iotservices/namespace-dropper", &isDeletion); err != nil {
		return reconcile.Result{}, err
	}
	if isDeletion {
		return reconcile.Result{}, nil
	}

	// for instance updates enforce a limit of 30 parallel update slots
	// generation=1 -> initial version
	// generation=2 -> finalizer, namespace and instanceId attached
	// generation>2 -> update
	// if an instance is already in phase Processing we have to continue that
	const parallelUpdateLimit = 30
	const requeueAfterSeconds = 30
	if instance.GetPhase() == cpv1alpha1.Available && instance.GetGeneration() > 2 {
		if currentUpdateCount >= parallelUpdateLimit {
			logger.Info("Throttling due to max. parallel updates. Requeuing request.", "ParallelUpdateLimit", parallelUpdateLimit, "RequeueAfterSeconds", requeueAfterSeconds)
			return reconcile.Result{RequeueAfter: requeueAfterSeconds * time.Second}, nil
		}
		currentUpdateCount++
	}

	// basic resources
	if err := r.createImagePullSecret(instance, secret.NewImagePullSecret(instance), operatorNamespace, &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if err := r.createConfigMap(instance, configmap.NewIoTServiceConfigMap(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if instance.GetEnableNetworkPolicies() {
		if _, err := r.createResource(instance, network_policy.NewDenyNotInstanceNetworkPolicy(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	}
	// deployments and statefulsets
	if err := r.createDeployment(instance, deployment.NewCockpitDeployment(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	// depending on parameter iotservices.core_spring.enabled either choose Core or CoreSpring+Processing
	if instance.GetCoreSpringEnabled() {
		if err := r.createDeployment(instance, deployment.NewCoreSpringDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		if err := r.createDeployment(instance, deployment.NewProcessingDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		// in case this is an update from OSGi to Spring, we need to delete the OSGi Core
		if err := r.deleteResource(instance, deployment.NewCoreDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	} else {
		if err := r.createDeployment(instance, deployment.NewCoreDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		// in case this is an update from Spring to OSGi, we need to delete the core-spring and processing
		if err := r.deleteResource(instance, deployment.NewCoreSpringDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		if err := r.deleteResource(instance, deployment.NewProcessingDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	}
	if err := r.createStatefulSet(instance, stateful_set.NewGatewayMqttStatefulSet(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if err := r.createDeployment(instance, deployment.NewGatewayRestDeployment(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if err := r.createDeployment(instance, deployment.NewHaproxyDeployment(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if err := r.createStatefulSet(instance, stateful_set.NewMessagingStatefulSet(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if instance.GetPostgresqlCommonSpec().Size == 1 {
		if err := r.createDeployment(instance, deployment.NewPostgresqlDeployment(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	}

	// services
	if _, err := r.createResource(instance, service.NewCockpitService(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	// depending on parameter iotservices.core_spring.enabled either choose Core or CoreSpring+Processing
	if instance.GetCoreSpringEnabled() {
		if _, err := r.createResource(instance, service.NewCoreSpringService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		if _, err := r.createResource(instance, service.NewProcessingService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		// in case this is an update from OSGi to Spring, we need to delete the OSGi Core
		if err := r.deleteResource(instance, service.NewCoreService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	} else {
		if _, err := r.createResource(instance, service.NewCoreService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		// in case this is an update from Spring to OSGi, we need to delete the core-spring and processing
		if err := r.deleteResource(instance, service.NewCoreSpringService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
		if err := r.deleteResource(instance, service.NewProcessingService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	}
	if _, err := r.createResource(instance, service.NewGatewayMqttService(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if _, err := r.createResource(instance, service.NewGatewayRestService(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if _, err := r.createResource(instance, service.NewHaproxyService(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if _, err := r.createResource(instance, service.NewMessagingService(instance), &requeue); err != nil {
		r.markAsFailed(instance, err)
		return reconcile.Result{}, err
	}
	if instance.GetPostgresqlCommonSpec().Size == 1 {
		if _, err := r.createResource(instance, service.NewPostgresqlService(instance), &requeue); err != nil {
			r.markAsFailed(instance, err)
			return reconcile.Result{}, err
		}
	}

	r.updateStatus(instance, &requeue)

	if !requeue && instance.GetPhase() == cpv1alpha1.Processing {
		return reconcile.Result{RequeueAfter: 5 * time.Second}, nil
	} else {
		return reconcile.Result{Requeue: requeue}, nil
	}
}

// createImagePullSecret checks if the image pull secret  exists in the instance namespace and if not,
// creates it as a copy of the image pull secret in the operator namespace
func (r *ReconcileIoTService) createImagePullSecret(i *latest.IoTService, s *corev1.Secret, operatorNamespace string, requeue *bool) error {
	logger := log.WithValues("InstanceNamespace", i.GetNamespace(), "InstanceId", i.GetInstanceId())
	// Check if image pull secret exists
	if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: s.Namespace, Name: s.Name}, &corev1.Secret{}); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Creating image pull secret in instance namespace")
			// retrieve the secret from the operator namespace
			imagePullSecretTemplate := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: operatorNamespace,
					Name:      s.Name,
				},
			}
			if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: imagePullSecretTemplate.Namespace, Name: imagePullSecretTemplate.Name}, imagePullSecretTemplate); err != nil {
				logger.Error(err, "Failed to get image pull secret from operator namespace")
				return err
			}
			// create a copy in the instance namespace
			s.Type = imagePullSecretTemplate.Type
			s.Data = imagePullSecretTemplate.Data
			// Set IoTService instance as the owner and controller
			if err := controllerutil.SetControllerReference(i.GetV1Object(), s, r.scheme); err != nil {
				logger.Error(err, "SetControllerReference for Secret failed")
				return err
			}

			if err := r.client.Create(context.TODO(), s); err != nil {
				logger.Error(err, "Failed to create image pull secret in instance namespace", "ImagePullSecret", s.Name)
				return err
			}
			*requeue = true
		} else {
			logger.Error(err, "Failed to get image pull secret from instance namespace")
			return err
		}
	}

	return nil
}

// createDeployment checks if the given deployment exists in the cluster and if not, creates it
// furthermore it aligns the state of the deployment regarding Docker image, replica count, resource specification and configuration
func (r *ReconcileIoTService) createDeployment(i *latest.IoTService, d *appsv1.Deployment, requeue *bool) error {
	logger := log.WithValues("InstanceNamespace", i.GetNamespace(), "InstanceId", i.GetInstanceId(), "Deployment.Name", d.Name)
	// make sure deployment is created in cluster
	if deploymentInCluster, err := r.createResource(i, d, requeue); err != nil {
		return err
	} else if !*requeue {
		// Ensure the deployment image, size and resources are the same as the spec
		var image string
		var size int32
		var resources corev1.ResourceRequirements
		switch d.Labels["component"] {
		case "cockpit":
			image = i.GetCockpitCommonSpec().Image
			size = i.GetCockpitCommonSpec().Size
			resources = i.GetCockpitCommonSpec().Resources
		case "core":
			image = i.GetCoreCommonSpec().Image
			size = i.GetCoreCommonSpec().Size
			resources = i.GetCoreCommonSpec().Resources
		case "core-spring":
			image = i.GetCoreSpringCommonSpec().Image
			size = i.GetCoreSpringCommonSpec().Size
			resources = i.GetCoreSpringCommonSpec().Resources
		case "processing":
			image = i.GetProcessingCommonSpec().Image
			size = i.GetProcessingCommonSpec().Size
			resources = i.GetProcessingCommonSpec().Resources
		case "gateway-rest":
			image = i.GetGatewayRestCommonSpec().Image
			size = i.GetGatewayRestCommonSpec().Size
			resources = i.GetGatewayRestCommonSpec().Resources
		case "haproxy":
			image = i.GetHaproxyCommonSpec().Image
			size = i.GetHaproxyCommonSpec().Size
			resources = i.GetHaproxyCommonSpec().Resources
		case "postgresql":
			image = i.GetPostgresqlCommonSpec().Image
			size = i.GetPostgresqlCommonSpec().Size
			resources = i.GetPostgresqlCommonSpec().Resources
		default:
			logger.Error(nil, "Component not found", "component", d.Labels["component"])
		}
		if reflect.DeepEqual(&resources, &deploymentInCluster.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Resources) != true {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(d)); err != nil {
				logger.Error(err, "Failed to update Resources")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Resources updated")
			*requeue = true
			return nil
		}
		if *deploymentInCluster.(*appsv1.Deployment).Spec.Replicas != size {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(d)); err != nil {
				logger.Error(err, "Failed to update Size")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Size changed", "size", size)
			*requeue = true
			return nil
		}
		if deploymentInCluster.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image != image {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(d)); err != nil {
				logger.Error(err, "Failed to update Image")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Image changed", "size", size)
			*requeue = true
			return nil
		}
		// Ensure the deployment has the current config
		if oldConfigChecksum, ok := deploymentInCluster.(*appsv1.Deployment).Spec.Template.ObjectMeta.Annotations["checksum/config"]; ok {
			newConfigChecksum := d.Spec.Template.ObjectMeta.Annotations["checksum/config"]
			if oldConfigChecksum != newConfigChecksum {
				if err := r.client.Update(context.TODO(), addUpdateTimestamp(d)); err != nil {
					logger.Error(err, "Failed to update Deployment")
				}
				logger.Info("New config checksum calculated.", "configChecksum", newConfigChecksum)
				*requeue = true
				return nil
			}
		}
	}
	return nil
}

// createStatefulSet checks if the given stateful set exists in the cluster and if not, creates it
// furthermore it aligns the state of the stateful set regarding Docker image, replica count, resource specification and configuration
func (r *ReconcileIoTService) createStatefulSet(i *latest.IoTService, set *appsv1.StatefulSet, requeue *bool) error {
	logger := log.WithValues("InstanceNamespace", i.GetNamespace(), "InstanceId", i.GetInstanceId(), "Set.Name", set.Name)
	// make sure stateful set is created in cluster
	if statefulSetInCluster, err := r.createResource(i, set, requeue); err != nil {
		return err
	} else if !*requeue {
		// Ensure the stateful set image, size and resources are the same as the spec
		var image string
		var size int32
		var resources corev1.ResourceRequirements
		switch set.Labels["component"] {
		case "gateway-mqtt":
			image = i.GetGatewayMqttCommonSpec().Image
			size = i.GetGatewayMqttCommonSpec().Size
			resources = i.GetGatewayMqttCommonSpec().Resources
		case "messaging":
			image = i.GetMmsCommonSpec().Image
			size = i.GetMmsCommonSpec().Size
			resources = i.GetMmsCommonSpec().Resources
		default:
			logger.Error(nil, "Component not found", "component", set.Labels["component"])
		}
		if reflect.DeepEqual(&resources, &statefulSetInCluster.(*appsv1.StatefulSet).Spec.Template.Spec.Containers[0].Resources) != true {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(set)); err != nil {
				logger.Error(err, "Failed to update Resources")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Resources updated")
			*requeue = true
			return nil
		}
		if *statefulSetInCluster.(*appsv1.StatefulSet).Spec.Replicas != size {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(set)); err != nil {
				logger.Error(err, "Failed to update Size")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Size changed", "size", size)
			*requeue = true
			return nil
		}
		if statefulSetInCluster.(*appsv1.StatefulSet).Spec.Template.Spec.Containers[0].Image != image {
			if err := r.client.Update(context.TODO(), addUpdateTimestamp(set)); err != nil {
				logger.Error(err, "Failed to update Image")
				return err
			}
			// Spec updated - return and requeue
			logger.Info("Image changed", "size", size)
			*requeue = true
			return nil
		}
		// Ensure the stateful set has the current config
		if oldConfigChecksum, ok := statefulSetInCluster.(*appsv1.StatefulSet).Spec.Template.ObjectMeta.Annotations["checksum/config"]; ok {
			newConfigChecksum := set.Spec.Template.ObjectMeta.Annotations["checksum/config"]
			if oldConfigChecksum != newConfigChecksum {
				if err := r.client.Update(context.TODO(), addUpdateTimestamp(set)); err != nil {
					logger.Error(err, "Failed to update StatefulSet")
				}
				logger.Info("New config checksum calculated.", "configChecksum", newConfigChecksum)
				*requeue = true
				return nil
			}
		}
	}
	return nil
}

// createConfigMap checks if the given config map exists in the cluster and if not, creates it
// furthermore it aligns the content of the config map with the spec of the IoTService object
func (r *ReconcileIoTService) createConfigMap(i *latest.IoTService, configMap *corev1.ConfigMap, requeue *bool) error {
	logger := log.WithValues("InstanceNamespace", i.GetNamespace(), "InstanceId", i.GetInstanceId(), "ConfigMap.Name", configMap.Name)
	// make sure config map is created in cluster
	if configMapInCluster, err := r.createResource(i, configMap, requeue); err != nil {
		return err
	} else if !*requeue {
		if reflect.DeepEqual(configMapInCluster.(*corev1.ConfigMap).Data, configMap.Data) != true {
			if err := r.client.Update(context.TODO(), configMap); err != nil {
				logger.Error(err, "Failed to update ConfigMap")
				return err
			}
			logger.Info("ConfigMap updated")
			*requeue = true
			return nil
		}
	}
	return nil
}

// createResource checks if the given resource exists in the cluster and if not, creates it
func (r *ReconcileIoTService) createResource(i *latest.IoTService, res KubernetesResource, requeue *bool) (KubernetesResource, error) {
	logger := log.WithValues("InstanceId", i.GetInstanceId(), "InstanceNamespace", i.GetNamespace(), "Resource.Kind", res.GetObjectKind().GroupVersionKind().Kind, "Resource.Name", res.GetName())
	// Check if resource exists
	var resourceInCluster KubernetesResource
	switch res.GetObjectKind().GroupVersionKind().Kind {
	case "Service":
		resourceInCluster = &corev1.Service{}
	case "NetworkPolicy":
		resourceInCluster = &networkingv1.NetworkPolicy{}
	case "ConfigMap":
		resourceInCluster = &corev1.ConfigMap{}
	case "Deployment":
		resourceInCluster = &appsv1.Deployment{}
	case "StatefulSet":
		resourceInCluster = &appsv1.StatefulSet{}
	}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: res.GetName(), Namespace: res.GetNamespace()}, resourceInCluster); err != nil {
		if errors.IsNotFound(err) {
			if err := controllerutil.SetControllerReference(i.GetV1Object(), res, r.scheme); err != nil {
				return nil, err
			}

			logger.Info("Creating a new resource")
			if err := r.client.Create(context.TODO(), res); err != nil {
				logger.Error(err, "Failed to create new resource")
				return nil, err
			}
			*requeue = true
			return resourceInCluster, nil
		} else {
			logger.Error(err, "Failed to get resource")
			return nil, err
		}
	}
	return resourceInCluster, nil
}

// deleteResource deletes the resource specified in the given template
func (r *ReconcileIoTService) deleteResource(i *latest.IoTService, resource KubernetesResource, requeue *bool) error {
	logger := log.WithValues("InstanceId", i.GetInstanceId(), "Resource.Kind", resource.GetObjectKind().GroupVersionKind().Kind, "Resource.Name", resource.GetName())
	if err := r.client.Delete(context.TODO(), resource); err != nil {
		if errors.IsNotFound(err) {
			// everything is fine, the resource did not exist at all
			return nil
		} else {
			logger.Error(err, "Failed to delete resource.")
			return err
		}
	}
	// resource was deleted
	*requeue = true
	logger.Info("Deleted resource.")
	return nil
}

// handleFinalizer takes care for:
// - adding the finalizer to the resource if not present
// - processing the finalizer if resources is marked for deletion
// - dropping the finalizer after successful processing
// params:
// - i - pointer to IoTService object
// - finalizerName - the name of the finalizer
// - isDeletion - OUTPUT parameter
func (r *ReconcileIoTService) handleFinalizer(i *latest.IoTService, finalizerName string, isDeletion *bool) error {
	logger := log.WithValues("InstanceId", i.GetInstanceId())
	// in case the IoT Service instance is in deletion, the DeletionTimestamp is set
	if i.GetDeletionTimestamp() != nil {
		*isDeletion = true
		if contains(i.GetFinalizers(), finalizerName) {
			// process the finalizer
			// TODO: in case we have more finalizers, extract the processing into dedicated methods
			switch finalizerName {
			case "com.sap.iotservices/namespace-dropper":
				// this finalizer deletes the namespace of the service instance explicitly.
				// although all resources (also the namespace) are marked as owned by the IoTService object,
				// there were random occasions, where automatic garbage collection did not work out. possible root cause:
				// the namespace-scoped custom resource IoTService owns the cluster-scoped resource Namespace
				// to avoid orphaned resources, this finalizer exists
				logger.Info("IoTService object is being deleted. Handling finalizer: deleting the instance namespace.")
				// delete namespace
				var dummy bool // we do not need requeue here, as this is the last action for the instance
				if err := r.deleteResource(i, namespace.NewNamespace(i), &dummy); err != nil {
					logger.Error(err, "Failed to delete instance namespace")
					return err
				}
			default:
				logger.Error(nil, "Finalizer processing not implemented!", "FinalizerName", finalizerName)
			}

			// drop finalizer and update object
			finalizers := remove(i.GetFinalizers(), finalizerName)
			i.SetFinalizers(finalizers)

			if err := r.client.Update(context.Background(), i.GetRuntimeObject()); err != nil {
				logger.Error(err, "Failed to remove finalizer")
				return err
			}

			logger.Info("Instance namespace deleted, finalizer removed.")
		} else {
			logger.Info("IoTService object is being deleted. Remaining resources are garbage collected.")
		}
		return nil
	} else {
		*isDeletion = false
		// instance is not being deleted, check if finalizer is already attached and do so if necessary
		if !contains(i.GetFinalizers(), finalizerName) {
			logger.Info("Attaching finalizer to IoTService object: ", finalizerName, i.GetFinalizers())

			finalizers := append(i.GetFinalizers(), finalizerName)
			i.SetFinalizers(finalizers)

			if err := r.client.Update(context.Background(), i.GetRuntimeObject()); err != nil {
				logger.Error(err, "Failed to attach finalizer")
				return err
			}
		}
		return nil
	}
}

// markAsFailed sets the phase of the IoTService object to "Failed" and sets the reason
// to the error message of the given error
func (r *ReconcileIoTService) markAsFailed(i *latest.IoTService, e error) {
	logger := log.WithValues("InstanceId", i.GetInstanceId())
	i.SetPhase(cpv1alpha1.Failed)
	i.SetReason(e.Error())
	if err := r.client.Update(context.TODO(), i.GetRuntimeObject()); err != nil {
		logger.Error(err, "Failed to update status: Failed", "reason", e.Error())
	}
}

// updateStatus updates the status of the IoTService object with
// - list of running pods
// - phase
func (r *ReconcileIoTService) updateStatus(i *latest.IoTService, requeue *bool) {
	logger := log.WithValues("InstanceId", i.GetInstanceId())
	// List the pods for this IoTService deployment
	podList := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(labelsForIoTService(i.GetName()))
	if err := r.client.List(context.TODO(), &client.ListOptions{Namespace: i.GetNamespace(), LabelSelector: labelSelector}, podList); err != nil {
		logger.Error(err, "Failed to list pods")
	}
	podNames := getPodNames(podList.Items)

	// Update the IoTService status with the pod names if needed
	if !reflect.DeepEqual(podNames, i.GetPods()) {
		i.SetPods(podNames)
		i.SetObservedGeneration(i.GetGeneration())
		if err := r.client.Status().Update(context.TODO(), i.GetRuntimeObject()); err != nil {
			logger.Error(err, "Failed to update status: Pods")
			*requeue = true
			return
		} else {
			logger.Info("Status updated: Pods")
		}
	}

	// if modifications have been triggered, requeue is set to true. in this case set the phase to "Processing"
	if *requeue && i.GetPhase() != cpv1alpha1.Processing {
		i.SetPhase(cpv1alpha1.Processing)
		i.SetObservedGeneration(i.GetGeneration())
		if err := r.client.Status().Update(context.TODO(), i.GetRuntimeObject()); err != nil {
			logger.Error(err, "Failed to update status: Processing")
			*requeue = true
			return
		} else {
			logger.Info("Status updated: Processing")
		}
	}

	// if no modifications have been triggered, start monitoring logic to check if everything is running
	if !*requeue {
		var reason string
		if status, err := r.checkInstanceAvailability(i, &reason); err != nil {
			logger.Error(err, "Failed to check instance availability")
		} else {
			i.SetPhase(status)
			i.SetObservedGeneration(i.GetGeneration())
			if status == cpv1alpha1.Failed {
				i.SetReason(reason)
			}
			if err := r.client.Status().Update(context.TODO(), i.GetRuntimeObject()); err != nil {
				logger.Error(err, fmt.Sprintf("Failed to update status: %v", status))
				*requeue = true
				return
			} else {
				logger.Info(fmt.Sprintf("Status updated: %v", status))
			}
		}
	}
}

// checkInstanceAvailability checks all necessary resources if they are available:
// - no unavailable pods in deployments
// - all statefulsets have required replica count and revision
// - public service has load balancer created
func (r *ReconcileIoTService) checkInstanceAvailability(i *latest.IoTService, reason *string) (cpv1alpha1.Phase, error) {
	logger := log.WithValues("InstanceId", i.GetInstanceId())
	// 1. check deployment status
	// list all deployments
	deploymentList := &appsv1.DeploymentList{}
	labelSelector := labels.SelectorFromSet(labelsForIoTService(i.GetName()))
	if err := r.client.List(context.TODO(), &client.ListOptions{Namespace: i.GetNamespace(), LabelSelector: labelSelector}, deploymentList); err != nil {
		logger.Error(err, "Failed to list deployments")
		return cpv1alpha1.Processing, err
	}
	requiredDeploymentNames := []string{
		"cockpit", "gateway-rest", "haproxy",
	}
	if i.GetCoreSpringEnabled() {
		requiredDeploymentNames = append(requiredDeploymentNames, "core-spring", "processing")
	} else {
		requiredDeploymentNames = append(requiredDeploymentNames, "core")
	}
	if i.GetPostgresqlCommonSpec().Size == 1 {
		requiredDeploymentNames = append(requiredDeploymentNames, "postgresql")
	}
	for _, d := range deploymentList.Items {
		componentName := d.Labels["component"]
		if contains(requiredDeploymentNames, componentName) {
			requiredDeploymentNames = remove(requiredDeploymentNames, componentName)
		}
	}
	if len(requiredDeploymentNames) != 0 {
		logger.Info("Not all deployments created yet.", "MissingDeployments", requiredDeploymentNames)
		return cpv1alpha1.Processing, nil
	}
	// all deployments created. check for each deployment if there are unavailable replicas
	for _, d := range deploymentList.Items {
		if d.Status.UnavailableReplicas > 0 {
			logger.Info("Some deployments still have unavailable replicas", "IncompleteDeployment", d.Name)
			// check if max progressing time has elapsed
			return cpv1alpha1.Processing, nil
		}
	}
	// all deployments available

	// 2. check statefulset status
	// list all statefulsets
	statefulSetList := &appsv1.StatefulSetList{}
	if err := r.client.List(context.TODO(), &client.ListOptions{Namespace: i.GetNamespace(), LabelSelector: labelSelector}, statefulSetList); err != nil {
		logger.Error(err, "Failed to list statefulsets")
		return cpv1alpha1.Processing, err
	}
	requiredStatefulSetNames := []string{
		"gateway-mqtt", "messaging",
	}
	for _, s := range statefulSetList.Items {
		componentName := s.Labels["component"]
		if contains(requiredStatefulSetNames, componentName) {
			requiredStatefulSetNames = remove(requiredStatefulSetNames, componentName)
		}
	}
	if len(requiredStatefulSetNames) != 0 {
		logger.Info("Not all statefulsets created yet", "MissingStatefulSets", requiredStatefulSetNames)
		return cpv1alpha1.Processing, nil
	}
	// all statefulsets created. check for each statefulset if it is up to the spec
	for _, s := range statefulSetList.Items {
		if *s.Spec.Replicas != s.Status.Replicas ||
			s.Status.Replicas != s.Status.ReadyReplicas ||
			s.Status.CurrentRevision != s.Status.UpdateRevision {
			logger.Info("Some statefulsets still have unavailable replicas", "IncompleteStatefulSet", s.Name)
			return cpv1alpha1.Processing, nil
		}
	}
	// all statefulsets available

	// 3. check public service
	// list public services
	serviceList := &corev1.ServiceList{}
	labelSelector = labels.SelectorFromSet(deployment.LabelsForHaproxy(i))
	if err := r.client.List(context.TODO(), &client.ListOptions{Namespace: i.GetNamespace(), LabelSelector: labelSelector}, serviceList); err != nil {
		logger.Error(err, "Failed to list public services")
		return cpv1alpha1.Processing, err
	}
	if len(serviceList.Items) == 0 {
		logger.Info("Public service not created yet")
		return cpv1alpha1.Processing, nil
	}
	// service created. check if load balancer already available
	if len(serviceList.Items[0].Status.LoadBalancer.Ingress) == 0 {
		logger.Info("Public service has not yet an attached load balancer")
		return cpv1alpha1.Processing, nil
	}
	// public service available

	currentUpdateCount--
	return cpv1alpha1.Available, nil
}
