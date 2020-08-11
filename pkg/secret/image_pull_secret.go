package secret

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/namespace"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewImagePullSecret (i *latest.IoTService) *corev1.Secret {
	const imagePullSecretName = "image-pull-secret"
	labels := namespace.LabelsForNamespace(i)
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: i.GetNamespace(),
			Name:      imagePullSecretName,
			Labels: labels,
		},
	}
}
