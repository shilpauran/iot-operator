package stateful_set

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/configmap"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NewGatewayMqttStatefulSet(wrapper *latest.IoTService) *appsv1.StatefulSet {

	cr := wrapper.GetLatest()

	annotations := map[string]string{
		"co.elastic.logs/module": "iot",
		"checksum/config":        configmap.ChecksumForIoTServiceConfigYAML(wrapper),
	}
	fsUserGroupId := int64(999)
	replicas := cr.Spec.Gateway.Mqtt.Size
	labels := LabelsForGatewayMqtt(wrapper)
	containerPort, _ := strconv.ParseInt(cr.Spec.Gateway.Mqtt.ContainerPort, 10, 32)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-mqtt-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "gateway-mqtt-" + cr.Name,
			Replicas:    &replicas,
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
						Image:           cr.Spec.Gateway.Mqtt.Image,
						Name:            "gateway-mqtt",
						ImagePullPolicy: corev1.PullPolicy(cr.Spec.ImagePullPolicy),
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(containerPort),
								Name:          "mqtt-tls",
							},
						},
						Env: []corev1.EnvVar{{
							Name:  "CONTAINER_PORT",
							Value: cr.Spec.Gateway.Mqtt.ContainerPort,
						}},
						Resources: cr.Spec.Gateway.Mqtt.Resources,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "config-volume",
								MountPath: "/etc/sap",
							},
							{
								Name:      "gateway-volume",
								MountPath: "/gateway/store",
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
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: &fsUserGroupId,
					},
					ImagePullSecrets: []corev1.LocalObjectReference{{
						Name: cr.Spec.ImagePullSecrets,
					}},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gateway-volume",
					Namespace: cr.Spec.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse(cr.Spec.Gateway.Mqtt.PersistentVolumeClaimSize),
						},
					},
				},
			}},
		},
	}
}

func LabelsForGatewayMqtt(cr *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  cr.GetName(),
		"component": "gateway-mqtt",
	}
}
