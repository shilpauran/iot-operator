package configmap

import (
	"crypto/md5"
	"encoding/hex"
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("configmap_iotservice")



func NewIoTServiceConfigMap(cr *latest.IoTService) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      NameForIoTService(cr),
			Namespace: cr.GetNamespace(),
			Labels:    LabelsForIoTService(cr),
		},
		Data: dataForIoTServiceConfigMap(cr),
	}
}

func LabelsForIoTService(cr *latest.IoTService) map[string]string {
	return map[string]string{
		"instance":  cr.GetName(),
		"component": "iotservice",
	}
}

func NameForIoTService(cr *latest.IoTService) string {
	return "iotservice-" + cr.GetName()
}

func ChecksumForIoTServiceConfigYAML(cr *latest.IoTService) string {
	hasher := md5.New()
	for k, v := range dataForIoTServiceConfigMap(cr) {
		if k == "config.yml" {
			hasher.Write([]byte(v))
			return hex.EncodeToString(hasher.Sum(nil))
		}
	}
	return ""
}

func dataForIoTServiceConfigMap(cr *latest.IoTService) map[string]string {
	log.WithValues("Request.Name", cr.GetName())
	bytes, err := cr.GetSerializedAppConfig()
	if err != nil {
		log.Error(err,"Failed to marshal configMap data")
	}
	data := map[string]string{
		"config.yml": string(bytes),
	}
	return data
}
