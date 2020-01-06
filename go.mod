module github.com/Azure/aad-pod-identity

go 1.12

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0
	github.com/Azure/azure-sdk-for-go v31.0.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.0
	github.com/Azure/go-autorest/autorest/adal v0.5.0
	github.com/Azure/go-autorest/autorest/azure/auth v0.1.0
	github.com/Azure/go-autorest/autorest/to v0.2.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.1.0 // indirect
	github.com/coreos/go-iptables v0.3.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7 // indirect
	github.com/google/go-cmp v0.3.0
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829 // indirect
	github.com/spf13/pflag v1.0.3
	go.opencensus.io v0.22.0
	golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
	k8s.io/klog v1.0.0
)
