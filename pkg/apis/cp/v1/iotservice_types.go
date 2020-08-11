package v1

import (
	cpv1alpha1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// The yaml tags are used only for the application ConfigMap

// IMPORTANT: Run "operator-sdk generate k8s" to regenerate code after modifying this file.
// Note that in operator-sdk version 0.4.0 there is a bug in the generation, which will cause the later step of
// "operator-sdk build ..." to fail. To work around this, you have to
// 1. comment out the declaration of IoTServiceCommonSpec.Resources and IoTServiceSpec.ProcessingServices
// 2. Run "operator-sdk generate k8s"
// 3. reactivate the commented-out lines
// 4. Run "operator-sdk build ..."

type AuditLogLastEventCache struct{
	ExpireAfterWriteSec int `json:"expireAfterWriteSec" yaml:"expireAfterWriteSec"`
	MaxSize             int `json:"maxSize" yaml:"maxSize"`
}

type AuditLog struct {
	Enabled      bool `json:"enabled" yaml:"enabled"`
	Domain       string `json:"domain" yaml:"domain"`
	User         string `json:"user" yaml:"user"`
	Password     string `json:"password" yaml:"password"`
	SubAccountId string `json:"subAccountId" yaml:"subAccountId"`
	AuditLogLastEventCache AuditLogLastEventCache `json:"auditLogLastEventCache" yaml:"auditLogLastEventCache"`
}

type Mms struct {
	cpv1alpha1.Mms `yaml:",inline"`
	ContainerPortJms            string `json:"container_port_jms" yaml:"container_port_jms"`
	ContainerPortJmsTls         string `json:"container_port_jms_tls" yaml:"container_port_jms_tls"`
	ContainerPortHazelCast      string `json:"container_port_hazelcast" yaml:"container_port_hazelcast"`
	PersistentVolumeClaimSize   string `json:"persistent_volume_claim_size" yaml:"persistent_volume_claim_size"`
}

type Core struct {
	cpv1alpha1.Core `yaml:",inline"`
	ContainerPortHttp           string `json:"container_port_http" yaml:"container_port_http"`
	ContainerPortHazelCast      string `json:"container_port_hazelcast" yaml:"container_port_hazelcast"`
}

type CoreSpring struct {
	cpv1alpha1.CoreSpring `yaml:",inline"`
	ContainerPort string `json:"container_port" yaml:"container_port"`
	PersistentVolumeClaimSize string `json:"persistent_volume_claim_size" yaml:"persistent_volume_claim_size"`
}

type ProcessingService struct {
	cpv1alpha1.ProcessingService `yaml:",inline"`
	ContainerPort string `json:"container_port" yaml:"container_port"`
	KafkaProducer struct {
		Acks               string `json:"acks" yaml:"acks"`
		MaxBlockInMs       string `json:"maxBlockInMs" yaml:"maxBlockInMs"`
		RequestTimeoutInMs string `json:"requestTimeoutInMs" yaml:"requestTimeoutInMs"`
		Retries            string `json:"retries" yaml:"retries"`
		RetryBackoffInMs   string `json:"retryBackoffInMs" yaml:"retryBackoffInMs"`
		MaxRequestSize     string `json:"maxRequestSize" yaml:"maxRequestSize"`
	} `json:"kafkaProducer" yaml:"kafkaProducer"`
	ThroughputConfiguration struct{
		SqlQueueCapacity           string `json:"sqlQueueCapacity" yaml:"sqlQueueCapacity"`
		HttpQueueCapacity          string `json:"httpQueueCapacity" yaml:"httpQueueCapacity"`
		KafkaQueueCapacity         string `json:"kafkaQueueCapacity" yaml:"kafkaQueueCapacity"`
		LeonardoIotQueueCapacity   string `json:"leonardoIotQueueCapacity" yaml:"leonardoIotQueueCapacity"`
		InternalSqlQueueCapacity   string `json:"internalSqlQueueCapacity" yaml:"internalSqlQueueCapacity"`
		Batching struct{
			BatchSize                               string `json:"batchSize" yaml:"batchSize"`
			ForceTimeoutMilliseconds                string `json:"forceTimeoutMilliseconds" yaml:"forceTimeoutMilliseconds"`
			CheckPeriodForCommitTimeoutMilliseconds string `json:"checkPeriodForCommitTimeoutMilliseconds" yaml:"checkPeriodForCommitTimeoutMilliseconds"`
			NumberOfThreads                         string `json:"numberOfThreads" yaml:"numberOfThreads"`
			MaxSizeMbytes                           string `json:"maxSizeMbytes" yaml:"maxSizeMbytes"`
			MaxThreadsWaitingForMeasures            string `json:"maxThreadsWaitingForMeasures" yaml:"maxThreadsWaitingForMeasures"`
			MaxThreadsProcessingMeasures            string `json:"maxThreadsProcessingMeasures" yaml:"maxThreadsProcessingMeasures"`
		} `json:"batching" yaml:"batching"`
	} `json:"throughputConfiguration" yaml:"throughputConfiguration"`
}

type Cockpit struct {
	cpv1alpha1.Cockpit `yaml:",inline"`
	ContainerPort    string `json:"container_port" yaml:"container_port"`
}

type Mqtt                    struct {
	cpv1alpha1.Mqtt `yaml:",inline"`
	AmqTempSizeMb              int    `json:"amq_temp_size_mb" yaml:"amq_temp_size_mb"`
	TenantExposed              bool    `json:"tenant_exposed" yaml:"tenant_exposed"`
	ContainerPort               string `json:"container_port" yaml:"container_port"`
	PersistentVolumeClaimSize   string `json:"persistent_volume_claim_size" yaml:"persistent_volume_claim_size"`
}

type Rest struct {
	cpv1alpha1.Rest `yaml:",inline"`
	CommandsTTL                 int `json:"commands_ttl" yaml:"commands_ttl"`
	ContainerPort               string `json:"container_port" yaml:"container_port"`
}

type Gateway struct {
	EnableDeviceVitality    bool `json:"enable_device_vitality" yaml:"enable_device_vitality"`
	EnableMeasureValidation bool `json:"enable_measure_validation" yaml:"enable_measure_validation"`
	Mqtt                    Mqtt `json:"mqtt" yaml:"mqtt"`
	Rest                    Rest `json:"rest" yaml:"rest"`
}

type Haproxy struct {
	cpv1alpha1.Haproxy `yaml:",inline"`
	ContainerPortHttp           string `json:"container_port_http" yaml:"container_port_http"`
	ContainerPortHttps          string `json:"container_port_https" yaml:"container_port_https"`
	ContainerPortMqtt           string `json:"container_port_mqtt" yaml:"container_port_mqtt"`
	ContainerPortJms            string `json:"container_port_jms" yaml:"container_port_jms"`
}

type Postgresql struct {
	cpv1alpha1.Postgresql `yaml:",inline"`
	ContainerPort               string `json:"container_port" yaml:"container_port"`
}

type MeasureMigration struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// IoTServiceSpec defines the desired state of IoTService
type IoTServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	Namespace                  string               `json:"instance_namespace" yaml:"instance_namespace"`
	InstanceId                 string               `json:"-" yaml:"instance_id"`
	InstanceOwnerPwd           string               `json:"instance_owner_pwd" yaml:"instance_owner_pwd"`
	ClientAuthCert             []string             `json:"client_auth_cert" yaml:"client_auth_cert"`
	ServerCert                 string               `json:"server_cert" yaml:"server_cert"`
	DomainName                 string               `json:"domain_name" yaml:"domain_name"`
	Mms                        Mms                  `json:"mms" yaml:"mms"`
	Core                       Core                 `json:"core" yaml:"core"`
	CoreSpring                 CoreSpring           `json:"core_spring" yaml:"core_spring"`
	DataMigration              cpv1alpha1.DataMigration        `json:"data_migration" yaml:"data_migration"`
	MeasureMigration           MeasureMigration     `json:"measure_migration" yaml:"measure_migration"`
	ProcessingService          ProcessingService    `json:"processing" yaml:"processing"`
	Cockpit                    Cockpit              `json:"cockpit" yaml:"cockpit"`
	Gateway                    Gateway              `json:"gateway" yaml:"gateway"`
	Haproxy                    Haproxy              `json:"haproxy" yaml:"haproxy"`
	Postgresql                 Postgresql           `json:"postgresql" yaml:"-"`
	Certificate                cpv1alpha1.Certificate          `json:"certificate" yaml:"certificate"`
	Metering                   cpv1alpha1.Metering             `json:"metering" yaml:"metering"`
	ProcessingServices         []cpv1alpha1.ProcessingServices `json:"processing_services" yaml:"processing_services"`
	Packetdb                   cpv1alpha1.Packetdb             `json:"packetdb" yaml:"packetdb"`
	LastAlive                  cpv1alpha1.LastAlive            `json:"last_alive" yaml:"last_alive"`
	EnableMultiProtocolSupport bool                 `json:"enable_multi_protocol_support" yaml:"enable_multi_protocol_support"`
	JwtValidationURL           string               `json:"jwt_validation_url" yaml:"jwt_validation_url"`
	ImagePullPolicy            string               `json:"imagePullPolicy" yaml:"-"`
	ImagePullSecrets           string               `json:"imagePullSecrets" yaml:"-"`
	EnableNetworkPolicies      bool                 `json:"enableNetworkPolicies" yaml:"-"`
	AuditLog                   AuditLog `json:"auditlog" yaml:"auditlog"`
	Onboarding                 string   `json:"on_boarding" yaml:"on_boarding"`
	EnableMonitoring           bool                 `json:"enable_monitoring" yaml:"-"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IoTService is the Schema for the iotservices API
// +k8s:openapi-gen=true
type IoTService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IoTServiceSpec   `json:"spec,omitempty"`
	Status cpv1alpha1.IoTServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IoTServiceList contains a list of IoTService
type IoTServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IoTService `json:"items"`
}

type AppConfig struct {
	Name       string
	IoTServices IoTServiceSpec
}

func init() {
	SchemeBuilder.Register(&IoTService{}, &IoTServiceList{})
}
