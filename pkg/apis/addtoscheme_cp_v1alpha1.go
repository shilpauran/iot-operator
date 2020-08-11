package apis

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme)
}
