package shared

import (
	cpv1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1"
	cpv1alpha1 "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/apis/cp/v1alpha1"
)

func ConvertToV1(src *cpv1alpha1.IoTService) *cpv1.IoTService {

	//logger := log.Log
	dst := new(cpv1.IoTService)

	// ObjectMeta
	dst.ObjectMeta = src.ObjectMeta

	// Status
	dst.Status = src.Status

	// Spec
	dst.Spec.Namespace = src.Spec.Namespace
	dst.Spec.InstanceId = src.Spec.InstanceId
	dst.Spec.InstanceOwnerPwd = src.Spec.InstanceOwnerPwd
	dst.Spec.ClientAuthCert = src.Spec.ClientAuthCert
	dst.Spec.ServerCert = src.Spec.ServerCert
	dst.Spec.DomainName = src.Spec.DomainName

	// Spec.MMS
	dst.Spec.Mms.ListeningAddress = src.Spec.Mms.ListeningAddress
	dst.Spec.Mms.MaxHeapSize = src.Spec.Mms.MaxHeapSize
	dst.Spec.Mms.AmqStoreSizeMb = src.Spec.Mms.AmqStoreSizeMb
	dst.Spec.Mms.IoTServiceCommonSpec = src.Spec.Mms.IoTServiceCommonSpec
	dst.Spec.Mms.ContainerPortHazelCast = "5701"
	dst.Spec.Mms.ContainerPortJms = "61619"
	dst.Spec.Mms.ContainerPortJmsTls = "61620"
	dst.Spec.Mms.PersistentVolumeClaimSize = "10Gi"

	// Spec.Core
	dst.Spec.Core.ListeningAddress = src.Spec.Core.ListeningAddress
	dst.Spec.Core.MaxHeapSize = src.Spec.Core.MaxHeapSize
	dst.Spec.Core.IoTServiceCommonSpec = src.Spec.Core.IoTServiceCommonSpec
	dst.Spec.Core.ContainerPortHttp = "8161"
	dst.Spec.Core.ContainerPortHazelCast = "5801"

	// Spec.CoreSpring
	dst.Spec.CoreSpring.Enabled = src.Spec.CoreSpring.Enabled
	dst.Spec.CoreSpring.MaxHeapSize = src.Spec.CoreSpring.MaxHeapSize
	dst.Spec.CoreSpring.Database = src.Spec.CoreSpring.Database
	dst.Spec.CoreSpring.IoTServiceCommonSpec = src.Spec.CoreSpring.IoTServiceCommonSpec
	dst.Spec.CoreSpring.ContainerPort = "8081"

	// Spec.DataMigration
	dst.Spec.DataMigration = src.Spec.DataMigration

	// Spec.ProcessingService
	dst.Spec.ProcessingService.MaxHeapSize = src.Spec.ProcessingService.MaxHeapSize
	dst.Spec.ProcessingService.Database = src.Spec.ProcessingService.Database
	dst.Spec.ProcessingService.Types = src.Spec.ProcessingService.Types
	dst.Spec.ProcessingService.IoTServiceCommonSpec = src.Spec.ProcessingService.IoTServiceCommonSpec
	dst.Spec.ProcessingService.ContainerPort = "8082"

	// Spec.Cockpit
	dst.Spec.Cockpit.MaxHeapSize = src.Spec.Cockpit.MaxHeapSize
	dst.Spec.Cockpit.ListeningAddress = src.Spec.Cockpit.ListeningAddress
	dst.Spec.Cockpit.IotConnect365 = src.Spec.Cockpit.IotConnect365
	dst.Spec.Cockpit.IoTServiceCommonSpec = src.Spec.Cockpit.IoTServiceCommonSpec
	dst.Spec.Cockpit.ContainerPort = "8080"

	// Spec.Gateway
	dst.Spec.Gateway.EnableDeviceVitality = src.Spec.Gateway.EnableDeviceVitality
	dst.Spec.Gateway.EnableMeasureValidation = src.Spec.Gateway.EnableMeasureValidation

	dst.Spec.Gateway.Mqtt.ListeningAddress = src.Spec.Gateway.Mqtt.ListeningAddress
	dst.Spec.Gateway.Mqtt.MaxHeapSize = src.Spec.Gateway.Mqtt.MaxHeapSize
	dst.Spec.Gateway.Mqtt.AmqStoreSizeMb = src.Spec.Gateway.Mqtt.AmqStoreSizeMb
	dst.Spec.Gateway.Mqtt.IoTServiceCommonSpec = src.Spec.Gateway.Mqtt.IoTServiceCommonSpec
	dst.Spec.Gateway.Mqtt.ContainerPort = "61628"
	dst.Spec.Gateway.Mqtt.PersistentVolumeClaimSize = "20Gi"

	dst.Spec.Gateway.Rest.ListeningAddress = src.Spec.Gateway.Rest.ListeningAddress
	dst.Spec.Gateway.Rest.MaxHeapSize = src.Spec.Gateway.Rest.MaxHeapSize
	dst.Spec.Gateway.Rest.IoTServiceCommonSpec = src.Spec.Gateway.Rest.IoTServiceCommonSpec
	dst.Spec.Gateway.Rest.ContainerPort = "8699"

	// Spec.Haproxy
	dst.Spec.Haproxy.IoTServiceCommonSpec = src.Spec.Haproxy.IoTServiceCommonSpec
	dst.Spec.Haproxy.ContainerPortHttp = "80"
	dst.Spec.Haproxy.ContainerPortHttps = "443"
	dst.Spec.Haproxy.ContainerPortMqtt = "8883"
	dst.Spec.Haproxy.ContainerPortJms = "61616"

	// Spec.Postgresql
	dst.Spec.Postgresql.IoTServiceCommonSpec = src.Spec.Postgresql.IoTServiceCommonSpec
	dst.Spec.Postgresql.ContainerPort = "5432"

	// Spec.Certificate
	dst.Spec.Certificate.Provider = src.Spec.Certificate.Provider
	dst.Spec.Certificate.RACertificatePassword = src.Spec.Certificate.RACertificatePassword
	dst.Spec.Certificate.TCS.CertificateAuthorityUrl = src.Spec.Certificate.TCS.CertificateAuthorityUrl

	// Spec.Metering
	dst.Spec.Metering = src.Spec.Metering

	// Spec.ProcessingServices
	dst.Spec.ProcessingServices = src.Spec.ProcessingServices

	// Spec.PacketDb
	dst.Spec.Packetdb = src.Spec.Packetdb

	// Spec.LastAlive
	dst.Spec.LastAlive = src.Spec.LastAlive

	dst.Spec.EnableMultiProtocolSupport = src.Spec.EnableMultiProtocolSupport
	dst.Spec.JwtValidationURL = src.Spec.JwtValidationURL
	dst.Spec.ImagePullPolicy = src.Spec.ImagePullPolicy
	dst.Spec.ImagePullSecrets = src.Spec.ImagePullSecrets
	dst.Spec.EnableNetworkPolicies = src.Spec.EnableNetworkPolicies

	// Status
	dst.Status.Pods = src.Status.Pods
	dst.Status.ObservedGeneration = src.Status.ObservedGeneration
	dst.Status.Phase = src.Status.Phase
	dst.Status.Reason = src.Status.Reason

	FillDefaultValues(dst)

	return dst
}


func FillDefaultValues(src *cpv1.IoTService) {

	// Spec.AuditLog
	if src.Spec.AuditLog.AuditLogLastEventCache.ExpireAfterWriteSec == 0 {
		src.Spec.AuditLog.AuditLogLastEventCache.ExpireAfterWriteSec = 60
	}

	// Spec.MMS
	if !(len(src.Spec.Mms.ContainerPortHazelCast) > 0) {
		src.Spec.Mms.ContainerPortHazelCast = "5701"
	}
	if !(len(src.Spec.Mms.ContainerPortJms) > 0) {
		src.Spec.Mms.ContainerPortJms = "61619"
	}
	if !(len(src.Spec.Mms.ContainerPortJmsTls) > 0) {
		src.Spec.Mms.ContainerPortJmsTls = "61620"
	}
	if !(len(src.Spec.Mms.PersistentVolumeClaimSize) > 0) {
		src.Spec.Mms.PersistentVolumeClaimSize = "10Gi"
	}

	// Spec.Core
	if !(len(src.Spec.Core.ContainerPortHttp) > 0) {
		src.Spec.Core.ContainerPortHttp = "8161"
	}
	if !(len(src.Spec.Core.ContainerPortHazelCast) > 0) {
		src.Spec.Core.ContainerPortHazelCast = "5801"
	}

	// Spec.CoreSpring
	if !(len(src.Spec.CoreSpring.ContainerPort) > 0) {
		src.Spec.CoreSpring.ContainerPort = "8081"
	}
	if !(len(src.Spec.CoreSpring.PersistentVolumeClaimSize) > 0) {
		src.Spec.CoreSpring.PersistentVolumeClaimSize = "2Gi"
	}

	// Spec.ProcessingService
	if !(len(src.Spec.ProcessingService.ContainerPort) > 0) {
		src.Spec.ProcessingService.ContainerPort = "8082"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.Acks) > 0) {
		src.Spec.ProcessingService.KafkaProducer.Acks = "1"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.MaxBlockInMs) > 0) {
		src.Spec.ProcessingService.KafkaProducer.MaxBlockInMs = "3000"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.RequestTimeoutInMs) > 0) {
		src.Spec.ProcessingService.KafkaProducer.RequestTimeoutInMs = "30000"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.Retries) > 0) {
		src.Spec.ProcessingService.KafkaProducer.Retries = "20"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.RetryBackoffInMs) > 0) {
		src.Spec.ProcessingService.KafkaProducer.RetryBackoffInMs = "5000"
	}
	if !(len(src.Spec.ProcessingService.KafkaProducer.MaxRequestSize) > 0) {
		src.Spec.ProcessingService.KafkaProducer.MaxRequestSize = "1048576"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.SqlQueueCapacity) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.SqlQueueCapacity = "200"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.HttpQueueCapacity) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.HttpQueueCapacity = "200"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.KafkaQueueCapacity) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.KafkaQueueCapacity = "200"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.LeonardoIotQueueCapacity) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.LeonardoIotQueueCapacity = "200"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.InternalSqlQueueCapacity) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.InternalSqlQueueCapacity = "200"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.BatchSize) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.BatchSize = "5000"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.ForceTimeoutMilliseconds) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.ForceTimeoutMilliseconds = "2000"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.CheckPeriodForCommitTimeoutMilliseconds) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.CheckPeriodForCommitTimeoutMilliseconds = "1000"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.NumberOfThreads) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.NumberOfThreads = "8"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxSizeMbytes) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxSizeMbytes = "128"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxThreadsWaitingForMeasures) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxThreadsWaitingForMeasures = "100"
	}
	if !(len(src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxThreadsProcessingMeasures) > 0) {
		src.Spec.ProcessingService.ThroughputConfiguration.Batching.MaxThreadsProcessingMeasures = "8"
	}

	// Spec.Cockpit
	if !(len(src.Spec.Cockpit.ContainerPort) > 0) {
		src.Spec.Cockpit.ContainerPort = "8080"
	}

	// Spec.Gateway
	if !(len(src.Spec.Gateway.Mqtt.ContainerPort) > 0) {
		src.Spec.Gateway.Mqtt.ContainerPort = "61628"
	}
	if src.Spec.Gateway.Mqtt.AmqTempSizeMb == 0 {
		src.Spec.Gateway.Mqtt.AmqTempSizeMb = 1000
	}
	if !(len(src.Spec.Gateway.Mqtt.PersistentVolumeClaimSize) > 0) {
		src.Spec.Gateway.Mqtt.PersistentVolumeClaimSize = "20Gi"
	}
	if !(len(src.Spec.Gateway.Rest.ContainerPort) > 0) {
		src.Spec.Gateway.Rest.ContainerPort = "8699"
	}
	if src.Spec.Gateway.Rest.CommandsTTL == 0 {
		src.Spec.Gateway.Rest.CommandsTTL = 21600000
	}

	// Always set tenantExposed to false, as this is currently the case.
	// This is required to avoid wrong configuration in case an instance starts
	// with wrong/missing configuration for this parameter
	src.Spec.Gateway.Mqtt.TenantExposed = false

	// Spec.Haproxy
	if !(len(src.Spec.Haproxy.ContainerPortHttp) > 0) {
		src.Spec.Haproxy.ContainerPortHttp = "80"
	}
	if !(len(src.Spec.Haproxy.ContainerPortHttps) > 0) {
		src.Spec.Haproxy.ContainerPortHttps = "443"
	}
	if !(len(src.Spec.Haproxy.ContainerPortMqtt) > 0) {
		src.Spec.Haproxy.ContainerPortMqtt = "8883"
	}
	if !(len(src.Spec.Haproxy.ContainerPortJms) > 0) {
		src.Spec.Haproxy.ContainerPortJms = "61616"
	}

	// Spec.Postgresql
	if !(len(src.Spec.Postgresql.ContainerPort) > 0) {
		src.Spec.Postgresql.ContainerPort = "5432"
	}

	// Spec.Onboarding
	if !(len(src.Spec.Onboarding) > 0) {
		src.Spec.Onboarding = "oneproduct"
	}

}
