package service

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/stateful_set"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewGatewayMqttService(wrapper *latest.IoTService) *corev1.Service {

	cr := wrapper.GetLatest()

	labels := stateful_set.LabelsForGatewayMqtt(wrapper)
	containerPort, _ := strconv.ParseInt(cr.Spec.Gateway.Mqtt.ContainerPort, 10, 32)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gateway-mqtt-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "mqtt-tls",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPort),
					TargetPort: intstr.FromInt(int(containerPort)),
				},
			},
		},
	}
}
