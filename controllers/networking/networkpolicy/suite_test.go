package networkpolicy

import (
	"context"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	//RunSpecsWithDefaultAndCustomReporters(t,
	//	"Controller Suite",
	//	[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = networkingv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	SetupWorkspaces([]Workspace{
		{Name: "foo", Namespaces: []Namespace{{Name: "foo1"}, {Name: "foo2"}}},
		{Name: "bar", Namespaces: []Namespace{{Name: "bar1"}, {Name: "bar2"}}},
	})

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

func SetupManager(ctx context.Context) {
	var stop chan struct{}

	BeforeEach(func() {
		stop = make(chan struct{})

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		err = (&Reconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
			Log:    ctrl.Log.WithName("controllers").WithName("NetworkPolicyHandler"),
		}).SetupWithManager(mgr)
		Expect(err).NotTo(HaveOccurred())

		go func() {
			defer GinkgoRecover()
			Expect(mgr.Start(stop)).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		close(stop)
	})
}

type Workspace struct {
	Name       string
	Namespaces []Namespace
}

type Namespace struct {
	Name string
}

var workspaces []Workspace

func SetupWorkspaces(w []Workspace) {
	workspaces = w

	for _, workspace := range workspaces {
		for _, namespace := range workspace.Namespaces {
			NewWorkspaceNamespaces(workspace.Name, namespace.Name)
		}
	}
}

func NewNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				constants.NamespaceLabelKey: name,
			},
		},
	}
}

func NewWorkspaceNamespaces(workspace string, namespaces ...string) {
	for _, namespace := range namespaces {
		ns := NewNamespace(namespace)
		ns.Labels[constants.WorkspaceLabelKey] = workspace

		err := k8sClient.Create(context.TODO(), ns)
		Expect(err).NotTo(HaveOccurred())
	}
}
