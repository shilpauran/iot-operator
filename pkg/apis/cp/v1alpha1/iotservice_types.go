package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
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

type IoTServiceCommonSpec struct {
	Image     string                      `json:"image" yaml:"-"`
	Size      int32                       `json:"size" yaml:"-"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty" yaml:"-"`
}

type Mms struct {
	ListeningAddress     string `json:"listening_address" yaml:"listening_address"`
	MaxHeapSize          string `json:"max_heap_size" yaml:"max_heap_size"`
	AmqStoreSizeMb       int    `json:"amq_store_size_mb" yaml:"amq_store_size_mb"`
	IoTServiceCommonSpec `yaml:"-"`
}

type Core struct {
	ListeningAddress     string `json:"listening_address" yaml:"listening_address"`
	MaxHeapSize          string `json:"max_heap_size" yaml:"max_heap_size"`
	IoTServiceCommonSpec `yaml:"-"`
}

type CoreSpring struct {
	Enabled     bool     `json:"enabled" yaml:"enabled"`
	MaxHeapSize string   `json:"max_heap_size" yaml:"max_heap_size"`
	Database    struct {
		Hostname        string `json:"hostname" yaml:"hostname"`
		Port            string `json:"port" yaml:"port"`
		Dbadminuser     string `json:"dbadminuser" yaml:"dbadminuser"`
		Dbadminpassword string `json:"dbadminpassword" yaml:"dbadminpassword"`
	} `json:"database" yaml:"database"`
	IoTServiceCommonSpec `yaml:"-"`
}

type DataMigration struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

type ProcessingService struct {
	MaxHeapSize string `json:"max_heap_size" yaml:"max_heap_size"`
	Database    struct {
		Hostname        string `json:"hostname" yaml:"hostname"`
		Port            string `json:"port" yaml:"port"`
		Dbadminuser     string `json:"dbadminuser" yaml:"dbadminuser"`
		Dbadminpassword string `json:"dbadminpassword" yaml:"dbadminpassword"`
		Dbname          string `json:"dbname" yaml:"dbname"`
	} `json:"database" yaml:"database"`
	Types struct {
		Kafka struct {
			Enabled bool `json:"enabled" yaml:"enabled"`
		} `json:"kafka" yaml:"kafka"`
		Iotae struct {
			Enabled      bool   `json:"enabled" yaml:"enabled"`
			UserName     string `json:"user" yaml:"user"`
			Password     string `json:"password" yaml:"password"`
			Brokers      string `json:"brokers" yaml:"brokers"`
			TokenUrl     string `json:"token_url" yaml:"token_url"`
			SubAccountId string `json:"subaccount_id" yaml:"subaccount_id"`
		} `json:"iotae" yaml:"iotae"`
		Sql struct {
			Enabled bool `json:"enabled" yaml:"enabled"`
		} `json:"sql" yaml:"sql"`
		Internalsql struct {
			Enabled bool `json:"enabled" yaml:"enabled"`
		} `json:"internalsql" yaml:"internalsql"`
		Http struct {
			Enabled bool `json:"enabled" yaml:"enabled"`
		} `json:"http" yaml:"http"`
	} `json:"types" yaml:"types"`
	IoTServiceCommonSpec `yaml:"-"`
}

type Cockpit struct {
	ListeningAddress string `json:"listening_address" yaml:"listening_address"`
	MaxHeapSize      string `json:"max_heap_size" yaml:"max_heap_size"`
	IotConnect365    struct {
		Enabled          bool   `json:"enabled" yaml:"enabled"`
		Host             string `json:"host" yaml:"host"`
		ApplicationToken string `json:"application_token" yaml:"application_token"`
	} `json:"iot_connect_365" yaml:"iot_connect_365"`
	IoTServiceCommonSpec `yaml:"-"`
}

type Mqtt                    struct {
	ListeningAddress     string `json:"listening_address" yaml:"listening_address"`
	MaxHeapSize          string `json:"max_heap_size" yaml:"max_heap_size"`
	AmqStoreSizeMb       int    `json:"amq_store_size_mb" yaml:"amq_store_size_mb"`
	IoTServiceCommonSpec `yaml:"-"`
}

type Rest struct {
	ListeningAddress     string `json:"listening_address" yaml:"listening_address"`
	MaxHeapSize          string `json:"max_heap_size" yaml:"max_heap_size"`
	IoTServiceCommonSpec `yaml:"-"`
}

type Gateway struct {
	EnableDeviceVitality    bool `json:"enable_device_vitality" yaml:"enable_device_vitality"`
	EnableMeasureValidation bool `json:"enable_measure_validation" yaml:"enable_measure_validation"`
	Mqtt Mqtt                         `json:"mqtt" yaml:"mqtt"`
	Rest Rest                         `json:"rest" yaml:"rest"`
}

type Haproxy struct {
	IoTServiceCommonSpec `yaml:"-"`
}

type Postgresql struct {
	IoTServiceCommonSpec `yaml:"-"`
}

type Certificate struct {
	Provider                string `json:"provider" yaml:"provider"`
	RACertificatePassword   string `json:"registration_authority_certificate_password" yaml:"registration_authority_certificate_password"`
	TCS struct {
		CertificateAuthorityUrl string `json:"certificate_authority_url,omitempty" yaml:"certificate_authority_url,omitempty"`
	} `json:"tcs,omitempty" yaml:"tcs,omitempty"`
}

type Metering struct {
	Maas struct {
		Enabled        bool   `json:"enabled" yaml:"enabled"`
		OrganizationID string `json:"organization_id" yaml:"organization_id"`
		SpaceID        string `json:"space_id" yaml:"space_id"`
		InstanceID     string `json:"instance_id" yaml:"instance_id"`
		Domain         string `json:"domain" yaml:"domain"`
		Region         string `json:"region" yaml:"region"`
		ClientID       string `json:"client_id" yaml:"client_id"`
		ClientSecret   string `json:"client_secret" yaml:"client_secret"`
	} `json:"maas" yaml:"maas"`
	Enabled        bool   `json:"enabled" yaml:"enabled"`
	OrganizationID string `json:"organization_id" yaml:"organization_id"`
	SpaceID        string `json:"space_id" yaml:"space_id"`
	Domain         string `json:"domain" yaml:"domain"`
	InstanceID     string `json:"instance_id" yaml:"instance_id"`
	ResourceID     string `json:"resource_id" yaml:"resource_id"`
	User           string `json:"user" yaml:"user"`
	Password       string `json:"password" yaml:"password"`
}

type ProcessingServices struct {
	Name       string `json:"name" yaml:"name"`
	Properties struct {
		ProcessingKafkaBrokers       string `json:"processing.kafka.brokers,omitempty" yaml:"processing.kafka.brokers,omitempty"`
		ProcessingKafkaTopicMetadata string `json:"processing.kafka.topic.metadata,omitempty" yaml:"processing.kafka.topic.metadata,omitempty"`
		ProcessingKafkaTopicMeasures string `json:"processing.kafka.topic.measures,omitempty" yaml:"processing.kafka.topic.measures,omitempty"`
		ProcessingKafkaUser          string `json:"processing.kafka.user,omitempty" yaml:"processing.kafka.user,omitempty"`
		ProcessingKafkaPassword      string `json:"processing.kafka.password,omitempty" yaml:"processing.kafka.password,omitempty"`
		ProcessingKafkaTokenurl      string `json:"processing.kafka.tokenurl,omitempty" yaml:"processing.kafka.tokenurl,omitempty"`
		ProcessingHttpUrl            string `json:"processing.http.url,omitempty" yaml:"processing.http.url,omitempty"`
		ProcessingHttpHeaders        []struct {
			HeaderName  string `json:"name" yaml:"name"`
			HeaderValue string `json:"value" yaml:"value"`
		} `json:"processing.http.headers" yaml:"processing.http.headers,omitempty"`
	} `json:"properties,omitempty" yaml:"properties,omitempty"`
}

type Packetdb struct {
	Hostname        string `json:"hostname" yaml:"hostname"`
	Port            string `json:"port" yaml:"port"`
	Dbadminuser     string `json:"dbadminuser" yaml:"dbadminuser"`
	Dbadminpassword string `json:"dbadminpassword" yaml:"dbadminpassword"`
}

type LastAlive struct {
	Enabled       bool `json:"enabled" yaml:"enabled"`
	PeriodSeconds int  `json:"period_seconds" yaml:"period_seconds"`
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
	DataMigration              DataMigration        `json:"data_migration" yaml:"data_migration"`
	ProcessingService          ProcessingService    `json:"processing" yaml:"processing"`
	Cockpit                    Cockpit              `json:"cockpit" yaml:"cockpit"`
	Gateway                    Gateway              `json:"gateway" yaml:"gateway"`
	Haproxy                    Haproxy              `json:"haproxy" yaml:"-"`
	Postgresql                 Postgresql           `json:"postgresql" yaml:"-"`
	Certificate                Certificate          `json:"certificate" yaml:"certificate"`
	Metering                   Metering             `json:"metering" yaml:"metering"`
	ProcessingServices         []ProcessingServices `json:"processing_services" yaml:"processing_services"`
	Packetdb                   Packetdb             `json:"packetdb" yaml:"packetdb"`
	LastAlive                  LastAlive            `json:"last_alive" yaml:"last_alive"`
	EnableMultiProtocolSupport bool                 `json:"enable_multi_protocol_support" yaml:"enable_multi_protocol_support"`
	JwtValidationURL           string               `json:"jwt_validation_url" yaml:"jwt_validation_url"`
	ImagePullPolicy            string               `json:"imagePullPolicy" yaml:"-"`
	ImagePullSecrets           string               `json:"imagePullSecrets" yaml:"-"`
	EnableNetworkPolicies      bool                 `json:"enableNetworkPolicies" yaml:"-"`
}

// Phase defines the current phase of processing of IoTService object
type Phase string
const (
	Processing Phase = "Processing"
	Available Phase = "Available"
	Failed Phase = "Failed"
)

// IoTServiceStatus defines the observed state of IoTService
type IoTServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Pods []string `json:"pods,omitempty" yaml:"pods,omitempty"`
	ObservedGeneration int64 `json:"observedGeneration,omitempty" yaml:"observedGeneration,omitempty"`
	Phase Phase `json:"phase" yaml:"phase"`
	Reason string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IoTService is the Schema for the iotservices API
// +k8s:openapi-gen=true
type IoTService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IoTServiceSpec   `json:"spec,omitempty"`
	Status IoTServiceStatus `json:"status,omitempty"`
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
