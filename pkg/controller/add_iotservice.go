package controller

import (
	"github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/controller/iotservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, iotservice.Add)
}
