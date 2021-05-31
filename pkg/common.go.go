package repository

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type AtlasObjects interface {
	v1.ObjectMetaAccessor
	schema.ObjectKind
}

func RepoOwnerRef(repo AtlasObjects) *v1.OwnerReference {
	return v1.NewControllerRef(repo.GetObjectMeta(), repo.GroupVersionKind())
}
