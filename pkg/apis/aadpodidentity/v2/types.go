package v2

import (
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventType int

const (
	PodCreated      EventType = 0
	PodDeleted      EventType = 1
	PodUpdated      EventType = 2
	IdentityCreated EventType = 3
	IdentityDeleted EventType = 4
	IdentityUpdated EventType = 5
	BindingCreated  EventType = 6
	BindingDeleted  EventType = 7
	BindingUpdated  EventType = 8
	Exit            EventType = 9
)

const (
	CRDGroup    = "aadpodidentity.k8s.io"
	CRDVersion  = "v2"
	CRDLabelKey = "aadpodidbinding"

	BehaviorKey = "aadpodidentity.k8s.io/Behavior"
	// BehaviorNamespaced ...
	BehaviorNamespaced = "namespaced"
	// AssignedIDCreated status indicates azure assigned identity is created
	AssignedIDCreated = "Created"
	// AssignedIDAssigned status indicates identity has been assigned to the node
	AssignedIDAssigned = "Assigned"
	// AssignedIDUnAssigned status indicates identity has been unassigned from the node
	AssignedIDUnAssigned = "Unassigned"
)

/*** Global data structures ***/

// AzureIdentity is the specification of the identity data structure.
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureIdentitySpec   `json:"Spec"`
	Status AzureIdentityStatus `json:"Status"`
}

// AzureIdentityBinding brings together the spec of matching pods and the identity which they can use.

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentityBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureIdentityBindingSpec   `json:"Spec"`
	Status AzureIdentityBindingStatus `json:"Status"`
}

//AzureAssignedIdentity contains the identity <-> pod mapping which is matched.

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureAssignedIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureAssignedIdentitySpec   `json:"Spec"`
	Status AzureAssignedIdentityStatus `json:"Status"`
}

//AzurePodIdentityException contains the pod selectors for all pods that don't require
// NMI to process and request token on their behalf.

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzurePodIdentityException struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzurePodIdentityExceptionSpec   `json:"Spec"`
	Status AzurePodIdentityExceptionStatus `json:"Status"`
}

/*** Lists ***/
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzureIdentity `json:"Items"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentityBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzureIdentityBinding `json:"Items"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureAssignedIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzureAssignedIdentity `json:"Items"`
}

//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzurePodIdentityExceptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AzurePodIdentityException `json:"Items"`
}

/*** AzureIdentity ***/
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IdentityType int

const (
	UserAssignedMSI  IdentityType = 0
	ServicePrincipal IdentityType = 1
)

type AzureIdentitySpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// UserAssignedMSI or Service Principal
	Type IdentityType `json:"Type"`

	// User assigned MSI resource id.
	ResourceID string `json:"ResourceID"`
	//Both User Assigned MSI and SP can use this field.
	ClientID string `json:"ClientID"`

	//Used for service principal
	ClientPassword api.SecretReference `json:"ClientPassword"`
	// Service principal tenant id.
	TenantID string `json:"TenantID"`
	// For service principal. Option param for specifying the  AD details.
	ADResourceID string `json:"ADResourceID"`
	ADEndpoint   string `json:"ADEndpoint"`

	Replicas *int32 `json:"Replicas"`
}

type AzureIdentityStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AvailableReplicas int32 `json:"AvailableReplicas"`
}

/*** AzureIdentityBinding ***/
type MatchType int

const (
	Explicit MatchType = 0
	Selector MatchType = 1
)

//AssignedIDState -  State indicator for the AssignedIdentity
type AssignedIDState int

const (
	//Created - Default state of the assigned identity
	Created AssignedIDState = 0
	//Assigned - When the underlying platform assignment of EMSI is complete
	//the state moves to assigned
	Assigned AssignedIDState = 1
)

const (
	AzureIDResource                = "azureidentities"
	AzureIDBindingResource         = "azureidentitybindings"
	AzureAssignedIDResource        = "azureassignedidentities"
	AzureIdentityExceptionResource = "azurepodidentityexceptions"
)

// AzureIdentityBindingSpec matches the pod with the Identity.
// Used to indicate the potential matches to look for between the pod/deployment
// and the identities present..
type AzureIdentityBindingSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AzureIdentity     string `json:"AzureIdentity"`
	Selector          string `json:"Selector"`
	// Weight is used to figure out which of the matching identities would be selected.
	Weight int `json:"Weight"`
}

type AzureIdentityBindingStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AvailableReplicas int32 `json:"AvailableReplicas"`
}

/*** AzureAssignedIdentitySpec ***/

//AzureAssignedIdentitySpec has the contents of Azure identity<->POD
type AzureAssignedIdentitySpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AzureIdentityRef  *AzureIdentity        `json:"AzureidentityRef"`
	AzureBindingRef   *AzureIdentityBinding `json:"AzurebindingRef"`
	Pod               string                `json:"Pod"`
	PodNamespace      string                `json:"PodNamespace"`
	NodeName          string                `json:"NodeName"`

	Replicas *int32 `json:"Replicas"`
}

// AzureAssignedIdentityStatus has the replica status of the resouce.
type AzureAssignedIdentityStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            string `json:"Status"`
	AvailableReplicas int32  `json:"AvailableReplicas"`
}

// AzurePodIdentityExceptionSpec matches pods with the selector defined.
// If request originates from a pod that matches the selector, nmi will
// proxy the request and send response back without any validation.
type AzurePodIdentityExceptionSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	PodLabels         map[string]string `json:"PodLabels"`
}

// AzurePodIdentityExceptionStatus ...
type AzurePodIdentityExceptionStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            string `json:"Status"`
}
