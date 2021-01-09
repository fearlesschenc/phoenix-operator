package initialize

import (
	"bytes"
	"context"
	"fmt"
	iamv1alpha2 "github.com/fearlesschenc/kubesphere/pkg/apis/iam/v1alpha2"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (init *initializer) isDevopsProject() bool {
	return init.obj.Labels[constants.DevOpsProjectLabelKey] != ""
}

func (init *initializer) getRoleTemplates() (*iamv1alpha2.RoleBaseList, error) {
	var label string
	if init.isDevopsProject() {
		label = fmt.Sprintf(iamv1alpha2.ScopeLabelFormat, iamv1alpha2.ScopeDevOps)
	} else {
		label = fmt.Sprintf(iamv1alpha2.ScopeLabelFormat, iamv1alpha2.ScopeNamespace)
	}

	templates := &iamv1alpha2.RoleBaseList{}
	if err := init.client.List(context.TODO(), templates, client.MatchingLabels(map[string]string{label: ""})); err != nil {
		return nil, err
	}

	return templates, nil
}

func (init *initializer) ensureDefaultRoleCreated() (reconcile.Result, error) {
	templates, err := init.getRoleTemplates()
	if err != nil {
		return reconcile.RequeueWithError(err)
	}

	for _, template := range templates.Items {
		desireRole := &rbacv1.Role{}
		if err = yaml.NewYAMLOrJSONDecoder(bytes.NewBuffer(template.Role.Raw), 1024).Decode(desireRole); err != nil {
			// in which error shouldn't happen here, we just skip it now
			continue
		}

		if desireRole.Kind != iamv1alpha2.ResourceKindRole {
			continue
		}

		currentRole := &rbacv1.Role{}
		if err := init.client.Get(context.TODO(), types.NamespacedName{Namespace: init.obj.Name, Name: desireRole.Name}, currentRole); err != nil {
			if !errors.IsNotFound(err) {
				return reconcile.RequeueWithError(err)
			}

			// create role
			desireRole.Namespace = init.obj.Name
			if err := init.client.Create(context.TODO(), desireRole); err != nil {
				return reconcile.RequeueWithError(err)
			}

			continue
		}

		roleModified := reconcile.ObjectUnchanged
		if !reflect.DeepEqual(desireRole.Labels, currentRole.Labels) {
			currentRole.Labels = desireRole.Labels
			roleModified = reconcile.ObjectChanged
		}
		if !reflect.DeepEqual(desireRole.Annotations, currentRole.Annotations) {
			currentRole.Annotations = desireRole.Annotations
			roleModified = reconcile.ObjectChanged
		}
		if !reflect.DeepEqual(desireRole.Rules, currentRole.Rules) {
			currentRole.Rules = desireRole.Rules
			roleModified = reconcile.ObjectChanged
		}
		if roleModified {
			if err := init.client.Update(context.TODO(), currentRole); err != nil {
				return reconcile.RequeueWithError(err)
			}
		}
	}

	return reconcile.Continue()
}

func (init *initializer) ensureCreatorAdminRole() (reconcile.Result, error) {
	creatorName := init.obj.Annotations[constants.CreatorAnnotationKey]
	creator := &iamv1alpha2.User{}
	if err := init.client.Get(context.TODO(), types.NamespacedName{Name: creatorName}, creator); err != nil {
		return reconcile.RequeueWithError(err)
	}

	if err := init.client.Create(context.TODO(), &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: init.obj.Name,
			Name:      fmt.Sprintf("%s-%s", creatorName, iamv1alpha2.NamespaceAdmin),
			Labels:    map[string]string{iamv1alpha2.UserReferenceLabel: creatorName},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     iamv1alpha2.ResourceKindRole,
			Name:     iamv1alpha2.NamespaceAdmin,
		},
		Subjects: []rbacv1.Subject{
			{
				Name:     creatorName,
				Kind:     iamv1alpha2.ResourceKindUser,
				APIGroup: rbacv1.GroupName,
			},
		},
	}); err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Continue()
}
