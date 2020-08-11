package deployment

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/configmap"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NewProcessingDeployment(wrapper *latest.IoTService) *appsv1.Deployment {

	cr := wrapper.GetLatest()

	annotations := map[string]string{
		"co.elastic.logs/module": "iot",
		"checksum/config":        configmap.ChecksumForIoTServiceConfigYAML(wrapper),
	}

	labels := LabelsForProcessing(wrapper)
	replicas := cr.Spec.ProcessingService.Size
	containerPort, _ := strconv.ParseInt(cr.Spec.ProcessingService.ContainerPort, 10, 32)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "processing-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           cr.Spec.ProcessingService.Image,
						Name:            "processing",
						ImagePullPolicy: corev1.PullPolicy(cr.Spec.ImagePullPolicy),
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(containerPort),
								Name:          "http",
							},
						},
						Env: []corev1.EnvVar{{
							Name:  "CONTAINER_PORT",
							Value: cr.Spec.ProcessingService.ContainerPort,
						}},
						Resources: cr.Spec.ProcessingService.Resources,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "config-volume",
								MountPath: "/etc/sap",
							},
						},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configmap.NameForIoTService(wrapper),
									},
								},
							},
						},
					},
					ImagePullSecrets: []corev1.LocalObjectReference{{
						Name: cr.Spec.ImagePullSecrets,
					}},
				},
			},
		},
	}
}

func LabelsForProcessing(cr *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  cr.GetName(),
		"component": "processing",
	}
}
