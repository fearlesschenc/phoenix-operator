package config

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fluxcd/source-controller/api/v1beta1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/*
cluster.yaml

workspaces:
- name: workspace1
  network_isolation: true
  hosts:
    - 127.0.0.1
	- 127.0.0.2

register application:
we don't need to provide lots of application definition as
they are defined in application repo.

applications:
- name: application1
  url: http://github.com/foo/bar.git
  secretRef:
	name: secret
*/

type WorkspaceSpec struct {
	Name string `json:"name" yaml:"name"`

	tenantv1alpha1.WorkspaceSpec
}

type ApplicationRepo struct {
	Name string `json:"name" yaml:"name"`

	v1beta1.GitRepositorySpec
}

type ClusterConfiguration struct {
	Workspaces       []WorkspaceSpec   `json:"workspaces" yaml:"workspaces"`
	ApplicationRepos []ApplicationRepo `json:"applications" yaml:"applications"`
}

func LoadClusterConfig(path string) (*ClusterConfiguration, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &ClusterConfiguration{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, err
	}

	return config, nil
}
