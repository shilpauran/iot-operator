package shared

import (
	cpv1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1"
	cpv1alpha1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"gopkg.in/yaml.v2"
)

type IIoTService interface {
	GetV1Object() v1.Object
	GetRuntimeObject() runtime.Object
	GetFinalizers() []string
	SetFinalizers([]string)
	GetInstanceId() string
	GetDeletionTimestamp() *v1.Time
	GetNamespace() string
	SetNamespace(string)
	GetName() string
	GetReason() string
	SetReason(string)
	GetPhase() cpv1alpha1.Phase
	SetPhase(cpv1alpha1.Phase)
	GetObservedGeneration() int64
	SetObservedGeneration(int64)
	GetPods() []string
	SetPods([]string)
	GetGeneration() int64
	SetGeneration(int64)
	GetGatewayMqttCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetGatewayRestCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetMmsCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetCoreSpringCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetCoreCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetProcessingCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetCockpitCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetHaproxyCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetPostgresqlCommonSpec() cpv1alpha1.IoTServiceCommonSpec
	GetEnableNetworkPolicies() bool
	GetCoreSpringEnabled() bool
}

type IoTService struct {
	V1 *cpv1.IoTService
	V1alpha1 *cpv1alpha1.IoTService
}

func (r IoTService) GetV1Object() v1.Object {
	if r.V1alpha1 != nil {
		return r.V1alpha1
	} else {
		return r.V1
	}
}

func (r IoTService) GetRuntimeObject() runtime.Object {
	if r.V1alpha1 != nil {
		return r.V1alpha1
	} else {
		return r.V1
	}
}

func (r IoTService) GetFinalizers() []string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Finalizers
	} else {
		return r.V1.Finalizers
	}
}

func (r IoTService) SetFinalizers(finalizers []string) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Finalizers = finalizers
	} else {
		r.V1.Finalizers = finalizers
	}
}

func (r IoTService) GetInstanceId() string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.InstanceId
	} else {
		return r.V1.Spec.InstanceId
	}
}

func (r IoTService) SetInstanceId(instanceId string) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Spec.InstanceId = instanceId
	} else {
		r.V1.Spec.InstanceId = instanceId
	}
}

func (r IoTService) GetDeletionTimestamp() *v1.Time {
	if r.V1alpha1 != nil {
		return r.V1alpha1.DeletionTimestamp
	} else {
		return r.V1.DeletionTimestamp
	}
}

func (r IoTService) GetNamespace() string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Namespace
	} else {
		return r.V1.Spec.Namespace
	}
}

func (r IoTService) SetNamespace(namespace string) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Spec.Namespace = namespace
	} else {
		r.V1.Spec.Namespace = namespace
	}
}

func (r IoTService) GetName() string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Name
	} else {
		return r.V1.Name
	}
}

func (r IoTService) GetReason() string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Status.Reason
	} else {
		return r.V1.Status.Reason
	}
}

func (r IoTService) SetReason(reason string) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Status.Reason = reason
	} else {
		r.V1.Status.Reason = reason
	}
}

func (r IoTService) GetPhase() cpv1alpha1.Phase {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Status.Phase
	} else {
		return r.V1.Status.Phase
	}
}

func (r IoTService) SetPhase(phase cpv1alpha1.Phase) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Status.Phase = phase
	} else {
		r.V1.Status.Phase = phase
	}
}

func (r IoTService) GetObservedGeneration() int64 {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Status.ObservedGeneration
	} else {
		return r.V1.Status.ObservedGeneration
	}
}

func (r IoTService) SetObservedGeneration(observedGeneration int64) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Status.ObservedGeneration = observedGeneration
	} else {
		r.V1.Status.ObservedGeneration = observedGeneration
	}
}

func (r IoTService) GetPods() []string {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Status.Pods
	} else {
		return r.V1.Status.Pods
	}
}

func (r IoTService) SetPods(pods []string) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Status.Pods = pods
	} else {
		r.V1.Status.Pods = pods
	}
}

func (r IoTService) GetGeneration() int64 {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Generation
	} else {
		return r.V1.Generation
	}
}

func (r IoTService) SetGeneration(generation int64) {
	if r.V1alpha1 != nil {
		r.V1alpha1.Generation = generation
	} else {
		r.V1.Generation = generation
	}
}

func (r IoTService) GetGatewayMqttCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Gateway.Mqtt.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Gateway.Mqtt.IoTServiceCommonSpec
	}
}

func (r IoTService) GetGatewayRestCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Gateway.Rest.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Gateway.Rest.IoTServiceCommonSpec
	}
}

func (r IoTService) GetMmsCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Mms.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Mms.IoTServiceCommonSpec
	}
}

func (r IoTService) GetCoreSpringCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.CoreSpring.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.CoreSpring.IoTServiceCommonSpec
	}
}

func (r IoTService) GetCoreCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Core.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Core.IoTServiceCommonSpec
	}
}

func (r IoTService) GetProcessingCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.ProcessingService.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.ProcessingService.IoTServiceCommonSpec
	}
}

func (r IoTService) GetCockpitCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Cockpit.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Cockpit.IoTServiceCommonSpec
	}
}

func (r IoTService) GetHaproxyCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Haproxy.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Haproxy.IoTServiceCommonSpec
	}
}

func (r IoTService) GetPostgresqlCommonSpec() cpv1alpha1.IoTServiceCommonSpec {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.Postgresql.IoTServiceCommonSpec
	} else {
		return r.V1.Spec.Postgresql.IoTServiceCommonSpec
	}
}

func (r IoTService) GetEnableNetworkPolicies() bool {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.EnableNetworkPolicies
	} else {
		return r.V1.Spec.EnableNetworkPolicies
	}
}

func (r IoTService) GetCoreSpringEnabled() bool {
	if r.V1alpha1 != nil {
		return r.V1alpha1.Spec.CoreSpring.Enabled
	} else {
		return r.V1.Spec.CoreSpring.Enabled
	}
}

func (r IoTService) GetLatest() *cpv1.IoTService {
	if r.V1alpha1 != nil {
		return ConvertToV1(r.V1alpha1)
	} else {
		return r.V1
	}
}

func (r IoTService) GetSerializedAppConfig() ([]byte, error) {
	if r.V1alpha1 != nil {
		// Alternative: always serialize the configMap with also v1 params
		v1 := ConvertToV1(r.V1alpha1)
		conf := cpv1.AppConfig{
			Name: v1.Name,
			IoTServices: v1.Spec,
		}
		return yaml.Marshal(conf)
	} else {
		conf := cpv1.AppConfig{
			Name: r.V1.Name,
			IoTServices: r.V1.Spec,
		}
		return yaml.Marshal(conf)
	}
}


