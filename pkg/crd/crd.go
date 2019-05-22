package crd

import (
	"fmt"
	"time"

	aadpodid "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity/v1"
	informers "github.com/Azure/aad-pod-identity/pkg/generated/informers/externalversions"
	listers "github.com/Azure/aad-pod-identity/pkg/generated/listers/aadpodidentity/v1"
	"github.com/Azure/aad-pod-identity/pkg/stats"

	clientset "github.com/Azure/aad-pod-identity/pkg/generated/clientset/versioned"
	"github.com/golang/glog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	rest             *rest.RESTClient
	BindingWatcher   cache.SharedInformer
	IdWatcher        cache.SharedInformer
	InformerFactory  informers.SharedInformerFactory
	IDLister         listers.AzureIdentityLister
	BindingLister    listers.AzureIdentityBindingLister
	AssignedIDLister listers.AzureAssignedIdentityLister
}

type ClientInt interface {
	Start(exit <-chan struct{})
	SyncCache(exit <-chan struct{})
	RemoveAssignedIdentity(assignedIdentity *aadpodid.AzureAssignedIdentity) error
	CreateAssignedIdentity(assignedIdentity *aadpodid.AzureAssignedIdentity) error
	ListBindings() (res []*aadpodid.AzureIdentityBinding, err error)
	ListAssignedIDs() (res []*aadpodid.AzureAssignedIdentity, err error)
	ListIds() (res []*aadpodid.AzureIdentity, err error)
	ListPodIds(podns, podname string) (*[]aadpodid.AzureIdentity, error)
}

func NewCRDClientLite(config *rest.Config) (crdClient *Client, err error) {

	informerClient, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	aadpodidentityinformers := informers.NewSharedInformerFactory(informerClient, 30*time.Second)
	assignedIDLister := aadpodidentityinformers.Aadpodidentity().V1().AzureAssignedIdentities().Lister()
	exit := make(chan struct{})
	go aadpodidentityinformers.Start(exit)

	return &Client{
		AssignedIDLister: assignedIDLister,
	}, nil
}

func NewCRDClient(config *rest.Config, eventCh chan aadpodid.EventType) (crdClient *Client, err error) {
	restClient, err := newRestClient(config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	bindingWatcher, err := newBindingWatcher(restClient, eventCh)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	idWatcher, err := newIdWatcher(restClient, eventCh)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	informerClient, err := clientset.NewForConfig(config)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	aadpodidentityinformers := informers.NewSharedInformerFactory(informerClient, 30*time.Second)

	idLister := aadpodidentityinformers.Aadpodidentity().V1().AzureIdentities().Lister()
	bindingLister := aadpodidentityinformers.Aadpodidentity().V1().AzureIdentityBindings().Lister()
	assignedIDLister := aadpodidentityinformers.Aadpodidentity().V1().AzureAssignedIdentities().Lister()

	return &Client{
		rest:             restClient,
		BindingWatcher:   bindingWatcher,
		IdWatcher:        idWatcher,
		InformerFactory:  aadpodidentityinformers,
		IDLister:         idLister,
		BindingLister:    bindingLister,
		AssignedIDLister: assignedIDLister,
	}, nil
}

func newRestClient(config *rest.Config) (r *rest.RESTClient, err error) {
	crdconfig := *config
	crdconfig.GroupVersion = &schema.GroupVersion{Group: aadpodid.CRDGroup, Version: aadpodid.CRDVersion}
	crdconfig.APIPath = "/apis"
	crdconfig.ContentType = runtime.ContentTypeJSON
	s := runtime.NewScheme()
	s.AddKnownTypes(*crdconfig.GroupVersion,
		&aadpodid.AzureIdentity{},
		&aadpodid.AzureIdentityList{},
		&aadpodid.AzureIdentityBinding{},
		&aadpodid.AzureIdentityBindingList{},
		&aadpodid.AzureAssignedIdentity{},
		&aadpodid.AzureAssignedIdentityList{})
	crdconfig.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(s)}

	//Client interacting with our CRDs
	restClient, err := rest.RESTClientFor(&crdconfig)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return restClient, nil
}

func newBindingWatcher(r *rest.RESTClient, eventCh chan aadpodid.EventType) (cache.SharedInformer, error) {
	azBindingWatcher := cache.NewSharedInformer(
		cache.NewListWatchFromClient(r, aadpodid.AzureIDBindingResource, v1.NamespaceAll, fields.Everything()),
		&aadpodid.AzureIdentityBinding{},
		time.Minute*10)
	if azBindingWatcher == nil {
		return nil, fmt.Errorf("Could not create watcher for %s", aadpodid.AzureIDBindingResource)
	}
	azBindingWatcher.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				glog.V(6).Infof("Binding created")
				eventCh <- aadpodid.BindingCreated
			},
			DeleteFunc: func(obj interface{}) {
				glog.V(6).Infof("Binding deleted")
				eventCh <- aadpodid.BindingDeleted
			},
			UpdateFunc: func(OldObj, newObj interface{}) {
				glog.V(6).Infof("Binding updated")
				eventCh <- aadpodid.BindingUpdated
			},
		},
	)
	return azBindingWatcher, nil
}

func newIdWatcher(r *rest.RESTClient, eventCh chan aadpodid.EventType) (cache.SharedInformer, error) {
	azIdWatcher := cache.NewSharedInformer(
		cache.NewListWatchFromClient(r, aadpodid.AzureIDResource, v1.NamespaceAll, fields.Everything()),
		&aadpodid.AzureIdentity{},
		time.Minute*10)
	if azIdWatcher == nil {
		return nil, fmt.Errorf("Could not create Identity watcher for %s", aadpodid.AzureIDResource)
	}
	azIdWatcher.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				glog.V(6).Infof("Identity created")
				eventCh <- aadpodid.IdentityCreated
			},
			DeleteFunc: func(obj interface{}) {
				glog.V(6).Infof("Identity deleted")
				eventCh <- aadpodid.IdentityDeleted
			},
			UpdateFunc: func(OldObj, newObj interface{}) {
				glog.V(6).Infof("Identity updated")
				eventCh <- aadpodid.IdentityUpdated
			},
		},
	)
	return azIdWatcher, nil
}

func (c *Client) Start(exit <-chan struct{}) {
	go c.BindingWatcher.Run(exit)
	go c.IdWatcher.Run(exit)
	go c.InformerFactory.Start(exit)
	glog.Info("CRD watchers started")
}

func (c *Client) SyncCache(exit <-chan struct{}) {
	if !cache.WaitForCacheSync(exit) {
		panic("Cache could not be synchronized")
	}
}

func (c *Client) RemoveAssignedIdentity(assignedIdentity *aadpodid.AzureAssignedIdentity) error {
	glog.V(6).Infof("Deletion of id named: %s", assignedIdentity.Name)
	begin := time.Now()
	err := c.rest.Delete().Namespace(assignedIdentity.Namespace).Resource("azureassignedidentities").Name(assignedIdentity.Name).Do().Error()
	stats.Update(stats.AssignedIDDel, time.Since(begin))
	return err
}

func (c *Client) CreateAssignedIdentity(assignedIdentity *aadpodid.AzureAssignedIdentity) error {
	glog.Infof("Got id %s to assign", assignedIdentity.Name)
	begin := time.Now()
	// Create a new AzureAssignedIdentity which maps the relationship between
	// id and pod
	glog.Infof("Creating assigned Id: %s", assignedIdentity.Name)
	var res aadpodid.AzureAssignedIdentity
	// TODO: Ensure that the status reflects the corresponding
	err := c.rest.Post().Namespace(assignedIdentity.Namespace).Resource("azureassignedidentities").Body(assignedIdentity).Do().Into(&res)
	if err != nil {
		glog.Error(err)
		return err
	}

	stats.Update(stats.AssignedIDAdd, time.Since(begin))
	//TODO: Update the status of the assign identity to indicate that the node assignment got done.
	return nil
}

func (c *Client) ListBindings() (res []*aadpodid.AzureIdentityBinding, err error) {
	begin := time.Now()
	res, err = c.BindingLister.List(labels.Everything())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	stats.Update(stats.BindingList, time.Since(begin))
	return res, nil
}

func (c *Client) ListAssignedIDs() (res []*aadpodid.AzureAssignedIdentity, err error) {
	begin := time.Now()
	res, err = c.AssignedIDLister.List(labels.Everything())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	stats.Update(stats.AssignedIDList, time.Since(begin))
	return res, nil
}

func (c *Client) ListIds() (res []*aadpodid.AzureIdentity, err error) {
	begin := time.Now()
	res, err = c.IDLister.List(labels.Everything())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	stats.Update(stats.IDList, time.Since(begin))
	return res, nil
}

//ListPodIds - given a pod with pod name space
func (c *Client) ListPodIds(podns, podname string) (*[]aadpodid.AzureIdentity, error) {

	azAssignedIDs, err := c.AssignedIDLister.List(labels.Everything())
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	var matchedIds []aadpodid.AzureIdentity
	for _, v := range azAssignedIDs {
		if v.Spec.Pod == podname && v.Spec.PodNamespace == podns {
			matchedIds = append(matchedIds, *v.Spec.AzureIdentityRef)
		}
	}

	return &matchedIds, nil
}
