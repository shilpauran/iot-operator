package service

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/stateful_set"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewMessagingService(wrapper *latest.IoTService) *corev1.Service {

	cr := wrapper.GetLatest()

	labels := stateful_set.LabelsForMessaging(wrapper)
	containerPortJms, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortJms, 10, 32)
	containerPortJmsTls, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortJmsTls, 10, 32)
	containerPortHazelCast, _ := strconv.ParseInt(cr.Spec.Mms.ContainerPortHazelCast, 10, 32)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "messaging-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "jms-no-tls",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortJms),
					TargetPort: intstr.FromInt(int(containerPortJms)),
				},
				{
					Name:       "jms-tls",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortJmsTls),
					TargetPort: intstr.FromInt(int(containerPortJmsTls)),
				},
				{
					Name:       "hazelcast",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortHazelCast),
					TargetPort: intstr.FromInt(int(containerPortHazelCast)),
				},
			},
		},
	}
}
