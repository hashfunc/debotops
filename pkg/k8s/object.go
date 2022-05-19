package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/pointer"
)

type Object interface {
	GetObjectMeta() metav1.Object
	GroupVersionKind() schema.GroupVersionKind
}

func NewOwnerReference(resource Object) *metav1.OwnerReference {
	meta := resource.GetObjectMeta()
	version, kind := resource.GroupVersionKind().ToAPIVersionAndKind()
	return &metav1.OwnerReference{
		APIVersion:         version,
		Kind:               kind,
		Name:               meta.GetName(),
		UID:                meta.GetUID(),
		BlockOwnerDeletion: pointer.Bool(true),
		Controller:         pointer.Bool(true),
	}
}

func NewOwnerReferences(resource Object) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*NewOwnerReference(resource),
	}
}
