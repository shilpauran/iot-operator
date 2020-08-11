package service

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/deployment"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewPostgresqlService(wrapper *latest.IoTService) *corev1.Service {

	cr := wrapper.GetLatest()

	labels := deployment.LabelsForPostgresql(wrapper)
	containerPort, _ := strconv.ParseInt(cr.Spec.Postgresql.ContainerPort, 10, 32)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "postgresql",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPort),
					TargetPort: intstr.FromInt(int(containerPort)),
				},
			},
		},
	}
}
