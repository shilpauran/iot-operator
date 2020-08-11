package namespace

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewNamespace (i *latest.IoTService) *corev1.Namespace {
	labels := LabelsForNamespace(i)
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: i.GetNamespace(),
			Labels: labels,
		},
	}
}

func LabelsForNamespace(i *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  i.GetName(),
	}
}
