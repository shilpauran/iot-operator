package iotservice

import (
	cpv1alpha1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_iotservice")

// Add creates a new IoTService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIoTService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("iotservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IoTService
	// we are running as cluster-scoped operator in order to be able to access other namespaces

	err = c.Watch(&source.Kind{Type: &cpv1alpha1.IoTService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Deployments and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource StatefulSets and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Services and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource ConfigMaps and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Secrets and requeue the owner IoTService
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cpv1alpha1.IoTService{},
	})
	if err != nil {
		return err
	}

	return nil
}

