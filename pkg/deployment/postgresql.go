package deployment

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NewPostgresqlDeployment(wrapper *latest.IoTService) *appsv1.Deployment {

	cr := wrapper.GetLatest()

	labels := LabelsForPostgresql(wrapper)
	replicas := cr.Spec.Postgresql.Size
	containerPort, _ := strconv.ParseInt(cr.Spec.Postgresql.ContainerPort, 10, 32)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-" + cr.Name,
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
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           cr.Spec.Postgresql.Image,
						Name:            "postgresql",
						ImagePullPolicy: corev1.PullPolicy(cr.Spec.ImagePullPolicy),
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: int32(containerPort),
								Name:          "postgresql",
							},
						},
						Env: []corev1.EnvVar{{
							Name:  "CONTAINER_PORT",
							Value: cr.Spec.Postgresql.ContainerPort,
						}},
						Resources: cr.Spec.Postgresql.Resources,
					}},
					ImagePullSecrets: []corev1.LocalObjectReference{{
						Name: cr.Spec.ImagePullSecrets,
					}},
				},
			},
		},
	}
}

func LabelsForPostgresql(cr *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  cr.GetName(),
		"component": "postgresql",
	}
}
