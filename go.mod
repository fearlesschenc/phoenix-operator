module github.com/fearlesschenc/phoenix-operator

go 1.13

require (
	github.com/fearlesschenc/kubesphere v0.1.0
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	k8s.io/api v0.19.3
	k8s.io/apiextensions-apiserver v0.19.2 // indirect
	k8s.io/apimachinery v0.19.3
	k8s.io/client-go v11.0.1-0.20190820062731-7e43eff7c80a+incompatible
	k8s.io/kubernetes v1.14.0
	sigs.k8s.io/controller-runtime v0.6.3
)

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0
	github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6
	github.com/ugorji/go => github.com/ugorji/go v0.0.0-20190128213124-ee1426cffec0
	go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20200520232829-54ba9589114f
	helm.sh/helm/v3 => github.com/openpitrix/helm/v3 v3.0.0-20200725015400-ebf6d7e5b2b0
	k8s.io/client-go => k8s.io/client-go v0.19.3
	k8s.io/kubernetes => k8s.io/kubernetes v1.13.0
)
