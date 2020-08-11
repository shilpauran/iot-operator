package e2e

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis"
	operator "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

/*
Pass in the CR List object IoTServiceList as an argument to AddToFrameworkScheme() because the framework needs to
ensure that the dynamic client has the REST mappings to query the API server for the CR type. The framework will
keep polling the API server for the mappings and timeout after 5 seconds, returning an error if the mappings were
not discovered in that time.
 */
func TestIoTService(t *testing.T) {
	IoTServiceList := &operator.IoTServiceList{
		TypeMeta: metav1.TypeMeta{
			Kind: "IoTService",
			APIVersion: "cp.iot.sap/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, IoTServiceList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
}