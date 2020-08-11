package network_policy

import (
	latest "github.wdf.sap.corp/iotservices/iotservices-k8s/iot-operator/pkg/shared"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Creates a network policy which applies to all pods of the specified instance, except the one running haproxy.
// The network policy allows ingress only from pods of this instance.
func NewDenyNotInstanceNetworkPolicy(cr *latest.IoTService) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NetworkPolicy",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deny-not-instance-" + cr.GetName(),
			Namespace: cr.GetNamespace(),
			Labels: map[string]string{
				"instance": cr.GetName(),
			},
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"instance": cr.GetName(),
				},
				MatchExpressions: []metav1.LabelSelectorRequirement{{
					Key:      "component",
					Operator: metav1.LabelSelectorOpNotIn,
					Values:   []string{"haproxy"},
				}},
			},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{
					From: []networkingv1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"instance": cr.GetName(),
								},
							},
						},
					},
				},
			},
		},
	}
}
