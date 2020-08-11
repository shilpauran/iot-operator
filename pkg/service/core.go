package service

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/deployment"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewCoreService(wrapper *latest.IoTService) *corev1.Service {

	cr := wrapper.GetLatest()

	labels := deployment.LabelsForCore(wrapper)
	containerPortHttp, _ := strconv.ParseInt(cr.Spec.Core.ContainerPortHttp, 10, 32)
	containerPortHazelCast, _ := strconv.ParseInt(cr.Spec.Core.ContainerPortHazelCast, 10, 32)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "core-" + cr.Name,
			Namespace: cr.Spec.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(containerPortHttp),
					TargetPort: intstr.FromInt(int(containerPortHttp)),
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
