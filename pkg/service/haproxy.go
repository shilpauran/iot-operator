package service

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/deployment"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewHaproxyService(wrapper *latest.IoTService) *corev1.Service {

	cr := wrapper.GetLatest()

	labels := deployment.LabelsForHaproxy(wrapper)
	containerPortHttp, _ := strconv.ParseInt(cr.Spec.Haproxy.ContainerPortHttp, 10, 32)
	containerPortHttps, _ := strconv.ParseInt(cr.Spec.Haproxy.ContainerPortHttps, 10, 32)
	containerPortMqtt, _ := strconv.ParseInt(cr.Spec.Haproxy.ContainerPortMqtt, 10, 32)
	containerPortJms, _ := strconv.ParseInt(cr.Spec.Haproxy.ContainerPortJms, 10, 32)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "iot-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector:              labels,
			Type:                  corev1.ServiceTypeLoadBalancer,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortHttp),
					TargetPort: intstr.FromInt(int(containerPortHttp)),
				},
				{
					Name:       "https",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortHttps),
					TargetPort: intstr.FromInt(int(containerPortHttps)),
				},
				{
					Name:       "mqtt-tls",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortMqtt),
					TargetPort: intstr.FromInt(int(containerPortMqtt)),
				},
				{
					Name:       "jms-tls",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortJms),
					TargetPort: intstr.FromInt(int(containerPortJms)),
				},
			},
		},
	}
}
