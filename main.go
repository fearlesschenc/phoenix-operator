/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	workloadv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/workload/v1alpha1"
	tenantcontroller "github.com/fearlesschenc/phoenix-operator/controllers/tenant"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(tenantv1alpha1.AddToScheme(scheme))
	utilruntime.Must(workloadv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "e6443dcd.phoenix.fearlesschenc.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// TODO
	// Setup with repo controller

	if err = (&tenantcontroller.WorkspaceReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Workspace"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Workspace")
		os.Exit(1)
	}
	if err = (&tenantcontroller.ClusterReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Cluster"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Cluster")
		os.Exit(1)
	}
	//if err = (&workloadcontroller.ApplicationReconciler{
	//	Client: mgr.GetClient(),
	//	Log:    ctrl.Log.WithName("controllers").WithName("Application"),
	//	Scheme: mgr.GetScheme(),
	//}).SetupWithManager(mgr); err != nil {
	//	setupLog.Error(err, "unable to create controller", "controller", "Application")
	//	os.Exit(1)
	//}
	//if err = (&workloadcontroller.OfflineApplicationReconciler{
	//	Client: mgr.GetClient(),
	//	Log:    ctrl.Log.WithName("controllers").WithName("OfflineApplication"),
	//	Scheme: mgr.GetScheme(),
	//}).SetupWithManager(mgr); err != nil {
	//	setupLog.Error(err, "unable to create controller", "controller", "OfflineApplication")
	//	os.Exit(1)
	//}

	// workspace webhook
	// hookServer := mgr.GetWebhookServer()
	// hookServer.Register("/mutate-tenant-v1alpha1-workspace", &webhook.Admission{Handler: &tenantv1alpha1.WorkspaceMutator{Client: mgr.GetClient()}})
	// hookServer.Register("/validate-tenant-v1alpha1-workspace", &webhook.Admission{Handler: &webhook2.WorkspaceValidator{Client: mgr.GetClient()}})

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
