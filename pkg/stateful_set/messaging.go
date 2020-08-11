package stateful_set

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/configmap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NewMessagingStatefulSet(wrapper *latest.IoTService) *appsv1.StatefulSet {

	cr := wrapper.GetLatest()

	annotations := map[string]string{
		"co.elastic.logs/module": "iot",
		"checksum/config":        configmap.ChecksumForIoTServiceConfigYAML(wrapper),
	}

	fsUserGroupId := int64(999)
	labels := LabelsForMessaging(wrapper)
	replicas := cr.Spec.Mms.Size
	containerPortJms, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortJms, 10, 32)
	containerPortJmsTls, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortJmsTls, 10, 32)
	containerPortHazelcast, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortHazelCast, 10, 32)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "messaging-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: "messaging-" + cr.Name,
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
						Image:           cr.Spec.Mms.Image,
						Name:            "messaging",
						ImagePullPolicy: corev1.PullPolicy(cr.Spec.ImagePullPolicy),
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(containerPortJms),
								Name:          "jms-no-tls",
							},
							{
								ContainerPort: int32(containerPortJmsTls),
								Name:          "jms-tls",
							},
							{
								ContainerPort: int32(containerPortHazelcast),
								Name:          "hazelcast",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "CONTAINER_PORT_JMS",
								Value: cr.Spec.Mms.ContainerPortJms,
							},
							{
								Name:  "CONTAINER_PORT_JMS_TLS",
								Value: cr.Spec.Mms.ContainerPortJmsTls,
							},
							{
								Name:  "CONTAINER_PORT_HAZELCAST",
								Value: cr.Spec.Mms.ContainerPortHazelCast,
							},
						},
						Resources: cr.Spec.Mms.Resources,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "config-volume",
								MountPath: "/etc/sap",
							},
							{
								Name:      "messaging-volume",
								MountPath: "/messaging/store",
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
					Name:      "messaging-volume",
					Namespace: cr.Spec.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{
						corev1.ReadWriteOnce,
					},
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							"storage": resource.MustParse(cr.Spec.Mms.PersistentVolumeClaimSize),
						},
					},
				},
			}},
		},
	}
}

func LabelsForMessaging(cr *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  cr.GetName(),
		"component": "messaging",
	}
}
