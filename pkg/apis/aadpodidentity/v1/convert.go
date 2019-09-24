package v1

import (
	internal "github.com/Azure/aad-pod-identity/pkg/apis/aadpodidentity/internal"
)

//AzureIdentityToInternal - converts the versioned AzureIdentity to internal representation
func AzureIdentityToInternal(from AzureIdentity) (to internal.AzureIdentity) {
	to.ResourceVersion = from.ResourceVersion
	to.Spec = AzureIDSpecToInternal(from.Spec)
	to.Status = AzureIDStatusToInternal(from.Status)
	return
}

// AzureIdentityBindingToInternal - converts the versioned AzureIdentityBinding to internal representation.
func AzureIdentityBindingToInternal(from AzureIdentityBinding) (to internal.AzureIdentityBinding) {
	to.ResourceVersion = from.ResourceVersion
	to.Spec = AzureIDBindingSpecToInternal(from.Spec)
	to.Status = AzureIDBindingStatusToInternal(from.Status)
	return
}

// AzureAssignedIdentityToInternal - converts the versioned AzureAssignedIdentity to internal representation
func AzureAssignedIdentityToInternal(from AzureAssignedIdentity) (to internal.AzureAssignedIdentity) {
	to.ResourceVersion = from.ResourceVersion
	to.Spec = AzureAssignedIDSpectToInternal(from.Spec)
	to.Status = AzureAssignedIDStatusToInternal(from.Status)
	return
}

// AzurePodIdentityExceptionToInternal - converts the versioned AzurePodIdentityException to internal representation.
func AzurePodIdentityExceptionToInternal(from AzurePodIdentityException) (to internal.AzurePodIdentityException) {
	to.ResourceVersion = from.ResourceVersion
	to.Spec = AzurePodIdentityExceptionSpecToInternal(from.Spec)
	to.Status = AzurePodIdentityExceptionStatusToInternal(from.Status)
	return
}

// AzureIDSpecToInternal - converts a versioned AzureIdentitySpec to internal one.
func AzureIDSpecToInternal(from AzureIdentitySpec) (to AzureIdentitySpec) {
	to.ObjectMeta = from.ObjectMeta
	to.Type = from.Type
	to.ResourceID = from.ResourceID
	to.ClientID = from.ClientID
	to.ClientPassword = from.ClientPassword
	to.TenantID = from.TenantID
	to.ADResourceID = from.ADResourceID
	to.ADEndpoint = from.ADEndpoint
	return
}

func AzureIDStatusToInternal(from AzureIdentityStatus) (to AzureIdentityStatus) {
	to.ObjectMeta = from.ObjectMeta
	return
}

func AzureIDBindingSpecToInternal(from AzureIdentityBindingSpec) (to AzureIdentityBindingSpec) {
	to.ObjectMeta = from.ObjectMeta
	to.AzureIdentity = from.AzureIdentity
	to.Selector = from.Selector
	to.Weight = from.Weight
	return
}



func AzureAssignedIDStatusToInternal(from AzureAssignedIdentityStatus) (to AzureAssignedIdentityStatus) {
	to.ObjectMeta = from.ObjectMeta
	to.Status = from.Status
	return
}

func