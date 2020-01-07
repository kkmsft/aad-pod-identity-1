package main

import (
	"flag"
	"time"

	"github.com/Azure/aad-pod-identity/pkg/mic"
	"github.com/Azure/aad-pod-identity/version"

	aadpodid "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity"
	"github.com/Azure/aad-pod-identity/pkg/crd"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var (
	kubeconfig        string
	cloudconfig       string
	forceNamespaced   bool
	versionInfo       bool
	syncRetryDuration time.Duration
	leaderElectionCfg mic.LeaderElectionConfig
	httpProbePort     string
)

func main() {
	defer klog.Flush()
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to the kube config")

	flag.Set("logtostderr", "true")
	flag.Set("v", "10")

	flag.Parse()

	if versionInfo {
		version.PrintVersionAndExit()
	}
	klog.Infof("Starting mic process. Version: %v. Build date: %v", version.MICVersion, version.BuildDate)
	if cloudconfig == "" {
		klog.Warningf("--cloudconfig not passed will use aadpodidentity-admin-secret")
	}
	if kubeconfig == "" {
		klog.Warningf("--kubeconfig not passed will use InClusterConfig")
	}

	klog.Infof("kubeconfig (%s) cloudconfig (%s)", kubeconfig, cloudconfig)
	config, err := buildConfig(kubeconfig)
	if err != nil {
		klog.Fatalf("Could not read config properly. Check the k8s config file, %+v", err)
	}

	eventCh := make(chan aadpodid.EventType, 100)
	crdClient, err := crd.NewCRDClient(config, eventCh)
	if err != nil {
		klog.Fatalf("%+v", err)
	}

	// Starts the leader election loop
	var exit <-chan struct{}
	crdClient.Start(exit)
	crdClient.SyncCache(exit)

	bindings, err := crdClient.ListBindings()
	if err != nil {
		klog.Fatalf("Could not get the bindings: %+v", err)
	}

	for _, v := range *bindings {
		//log.Infof("\n=========>")
		klog.Infof("\n%s", v.Spec.Selector)
		//log.Infof("\n<=========>")
	}

	assignedIDs, err := crdClient.ListAssignedIDs()
	if err != nil {
		klog.Fatalf("Could not get assigned ID")
	}

	for _, a := range *assignedIDs {
		klog.Infof("\n%s\n", a.Status.Status)
	}

	klog.Info("Done !")
}

// Create the client config. Use kubeconfig if given, otherwise assume in-cluster.
func buildConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	return rest.InClusterConfig()
}
