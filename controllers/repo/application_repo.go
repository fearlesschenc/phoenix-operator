package repo

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/pkg/config"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
)

const ApplicationConfiguration = "application.yaml"

func isApplicationRepo(rootFiles []os.FileInfo) bool {
	for _, file := range rootFiles {
		if file.Name() == ApplicationConfiguration {
			return true
		}
	}

	return false
}

// updateApplication map git repo changes to Application
func (r *Reconciler) updateApplication(ctx context.Context, repo *sourcev1beta1.GitRepository, dir string) (ctrl.Result, error) {
	log := r.Log.WithValues("type", "cluster", "name", repo.Name)

	appConfig, err := config.LoadApplicationConfig(filepath.Join(dir, ApplicationConfiguration))
	if err != nil {
		return ctrl.Result{}, err
	}

	// TODO
	log.Info(appConfig.Services[0])

	return ctrl.Result{}, nil
}
