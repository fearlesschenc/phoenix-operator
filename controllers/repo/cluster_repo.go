package repo

//const (
//	ClusterConfiguration = "cluster.yaml"
//	DefaultClusterName   = "default"
//)
//
//func isClusterRepo(rootFiles []os.FileInfo) bool {
//	for _, file := range rootFiles {
//		if file.NodeName() == ClusterConfiguration {
//			return true
//		}
//	}
//
//	return false
//}

//func (r *Reconciler) updateCluster(ctx context.Context, repo *sourcev1beta1.GitRepository, artifactDir string) (ctrl.Result, error) {
//	log := r.Log.WithValues("type", "cluster", "name", repo.NodeName)

//created := true
//cluster := &tenantv1alpha1.Cluster{}
//if err := r.Get(ctx, types.NamespacedName{NodeName: DefaultClusterName}, cluster); err != nil {
//	if !errors.IsNotFound(err) {
//		return ctrl.Result{}, err
//	}

//	created = false
//	// init cluster
//	cluster.NodeName = DefaultClusterName
//}

//// map cluster repo fact to cluster spec
//clusterConfig, err := config.LoadClusterConfig(filepath.Join(artifactDir, ClusterConfiguration))
//if err != nil {
//	log.Error(err, "unable to read cluster configuration")
//	return ctrl.Result{}, err
//}

//// empty workspace list and reassign
//cluster.Spec.Workspaces = []tenantv1alpha1.WorkspaceTemplate{}
//for _, workspace := range clusterConfig.Workspaces {
//	cluster.Spec.Workspaces = append(cluster.Spec.Workspaces, tenantv1alpha1.WorkspaceTemplate{
//		ObjectMeta: metav1.ObjectMeta{NodeName: workspace.NodeName},
//		Spec: tenantv1alpha1.WorkspaceSpec{
//			NetworkIsolation: workspace.NetworkIsolation,
//			Hosts:                   workspace.Hosts,
//		},
//	})
//}

////cluster.Spec.Applications = []sourcev1beta1.GitRepositorySpec{}
////for _, repo := range clusterConfig.ApplicationRepos {
////	cluster.Spec.Applications = append(cluster.Spec.Applications, sourcev1beta1.GitRepositorySpec{
////		URL:          repo.URL,
////		SecretRef:    repo.SecretRef,
////		Interval:     repo.Interval,
////		Timeout:      repo.Timeout,
////		Reference:    repo.Reference,
////		Verification: repo.Verification,
////		Ignore:       repo.Ignore,
////		Suspend:      repo.Suspend,
////	})
////}

//if !created {
//	if err := r.Create(ctx, cluster); err != nil {
//		return ctrl.Result{}, err
//	}
//} else {
//	if err := r.Update(ctx, cluster); err != nil {
//		return ctrl.Result{}, err
//	}
//}

//	return ctrl.Result{}, nil
//}
