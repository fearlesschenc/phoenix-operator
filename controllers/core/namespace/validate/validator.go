package validate

import (
	"context"
	iamv1alpha2 "github.com/fearlesschenc/kubesphere/pkg/apis/iam/v1alpha2"
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type validator struct {
	client client.Client
	logger logr.Logger
	obj    *corev1.Namespace
}

func (v *validator) ensureWorkspaceValid() (reconcile.Result, error) {
	workspaceName := v.obj.Labels[constants.WorkspaceLabelKey]
	if workspaceName == "" {
		v.logger.Info("empty workspace name")
		return reconcile.Stop()
	}

	workspace := &tenantv1alpha1.Workspace{}
	if err := v.client.Get(context.TODO(), types.NamespacedName{Name: workspaceName}, workspace); err != nil {
		if errors.IsNotFound(err) {
			v.logger.Info("workspace not found")
		} else {
			v.logger.Error(err, "get workspace error")
		}

		return reconcile.Stop()
	}

	return reconcile.Continue()
}

func (v *validator) ensureCreatorValid() (reconcile.Result, error) {
	creatorName := v.obj.Annotations[constants.CreatorAnnotationKey]
	if creatorName == "" {
		v.logger.Info("empty creator name")
		return reconcile.Stop()
	}

	creator := &iamv1alpha2.User{}
	if err := v.client.Get(context.TODO(), types.NamespacedName{Name: creatorName}, creator); err != nil {
		return reconcile.RequeueOnErrorOrStop(client.IgnoreNotFound(err))
	}

	return reconcile.Continue()
}

func (v *validator) EnsureValidated() (reconcile.Result, error) {
	return reconcile.RunSubRoutine([]reconcile.SubroutineFunc{
		v.ensureWorkspaceValid,
		v.ensureCreatorValid,
	})
}
