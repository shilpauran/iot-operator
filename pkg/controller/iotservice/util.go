package iotservice

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

// generic interface which all Kubernetes Resources have to implement.
type KubernetesResource interface {
	metav1.Object
	runtime.Object
}

// labelsForIoTService returns the labels for selecting the resources
// belonging to the given IoTService CR name.
func labelsForIoTService(name string) map[string]string {
	return map[string]string{"instance": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		if pod.Status.Phase == corev1.PodRunning {
			podNames = append(podNames, pod.Name)
		}
	}
	return podNames
}

// addUpdateTimestamp adds an annotation "updateTimestamp" to the resource and returns it for chained calls
func addUpdateTimestamp(r KubernetesResource) KubernetesResource {
	if r.GetAnnotations() == nil {
		r.SetAnnotations(map[string]string{})
	}
	r.GetAnnotations()["updateTimestamp"] = time.Now().UTC().Format(time.RFC3339)
	return r
}

// contains tells whether a contains x.
func contains(a []string, s string) bool {
	for _, n := range a {
		if s == n {
			return true
		}
	}
	return false
}

// remove removes all strings from the array which are equal to s
func remove(a []string, s string) (result []string) {
	for _, n := range a {
		if s == n {
			continue
		}
		result = append(result, n)
	}
	return result
}
