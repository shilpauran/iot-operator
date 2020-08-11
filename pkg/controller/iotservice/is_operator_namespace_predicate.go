package iotservice

import (
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

var _ predicate.Predicate = &isOperatorNamespacePredicate{}

// that sigs.k8s.io/controller-runtime/pkg/predicate.Predicate only accepts requests where the namespace
// starts with the operator namespace
type isOperatorNamespacePredicate struct {
	operatorNamespace string
}

func (t isOperatorNamespacePredicate) Create(e event.CreateEvent) bool {
	return strings.HasPrefix(e.Meta.GetNamespace(), t.operatorNamespace)
}

func (t isOperatorNamespacePredicate) Delete(e event.DeleteEvent) bool {
	return strings.HasPrefix(e.Meta.GetNamespace(), t.operatorNamespace)
}

func (t isOperatorNamespacePredicate) Update(e event.UpdateEvent) bool {
	return strings.HasPrefix(e.MetaOld.GetNamespace(), t.operatorNamespace)
}

func (t isOperatorNamespacePredicate) Generic(e event.GenericEvent) bool {
	return strings.HasPrefix(e.Meta.GetNamespace(), t.operatorNamespace)
}

