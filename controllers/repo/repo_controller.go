package repo

import (
	"context"
	"fmt"
	"github.com/fluxcd/pkg/untar"
	sourcev1beta1 "github.com/fluxcd/source-controller/api/v1beta1"
	"github.com/go-logr/logr"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile dump GitRepository to corresponding CR based on repository
// structure.
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	repo := &sourcev1beta1.GitRepository{}
	if err := r.Get(ctx, req.NamespacedName, repo); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log := r.Log.WithValues(strings.ToLower(repo.Kind), req.NamespacedName)
	log.Info("New revision detected", "revision", repo.Status.Artifact.Revision)

	tmpDir, err := ioutil.TempDir("", repo.Name)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create temp dir, error: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	summary, err := r.fetchArtifact(ctx, *repo, tmpDir)
	if err != nil {
		log.Error(err, "unable to fetch artifact")
		return ctrl.Result{}, err
	}
	log.Info(summary)

	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list files, error: %w", err)
	}

	if !isApplicationRepo(files) {
		// ignore invalid
		log.Info("unknown repo type")
		return ctrl.Result{}, nil
	}

	return r.updateApplication(ctx, repo, tmpDir)
}

func (r *Reconciler) fetchArtifact(ctx context.Context, repo sourcev1beta1.GitRepository, dir string) (string, error) {
	if repo.Status.Artifact == nil {
		return "", fmt.Errorf("repository %s does not containt an artifact", repo.Name)
	}

	url := repo.Status.Artifact.URL
	if hostname := os.Getenv("SOURCE_HOST"); hostname != "" {
		url = fmt.Sprintf("http://%s/gitrepository/%s/%s/latest.tar.gz", hostname, repo.Namespace, repo.Name)
	}

	// download the tarball
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request, error : %w", err)
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("failed to download artifact from %s, error: %w", url, err)
	}
	defer resp.Body.Close()

	// check response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download artifact, status: %s", resp.Status)
	}

	// extract
	summary, err := untar.Untar(resp.Body, dir)
	if err != nil {
		return "", fmt.Errorf("failed to untar artifact, error: %w", err)
	}

	return summary, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sourcev1beta1.GitRepository{}).
		Complete(r)
}
