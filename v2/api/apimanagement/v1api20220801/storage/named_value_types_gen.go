// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package storage

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +kubebuilder:rbac:groups=apimanagement.azure.com,resources=namedvalues,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apimanagement.azure.com,resources={namedvalues/status,namedvalues/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20220801.NamedValue
// Generator information:
// - Generated from: /apimanagement/resource-manager/Microsoft.ApiManagement/stable/2022-08-01/apimnamedvalues.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/namedValues/{namedValueId}
type NamedValue struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Service_NamedValue_Spec   `json:"spec,omitempty"`
	Status            Service_NamedValue_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &NamedValue{}

// GetConditions returns the conditions of the resource
func (value *NamedValue) GetConditions() conditions.Conditions {
	return value.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (value *NamedValue) SetConditions(conditions conditions.Conditions) {
	value.Status.Conditions = conditions
}

var _ genruntime.KubernetesResource = &NamedValue{}

// AzureName returns the Azure name of the resource
func (value *NamedValue) AzureName() string {
	return value.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2022-08-01"
func (value NamedValue) GetAPIVersion() string {
	return string(APIVersion_Value)
}

// GetResourceScope returns the scope of the resource
func (value *NamedValue) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (value *NamedValue) GetSpec() genruntime.ConvertibleSpec {
	return &value.Spec
}

// GetStatus returns the status of this resource
func (value *NamedValue) GetStatus() genruntime.ConvertibleStatus {
	return &value.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (value *NamedValue) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationDelete,
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationHead,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.ApiManagement/service/namedValues"
func (value *NamedValue) GetType() string {
	return "Microsoft.ApiManagement/service/namedValues"
}

// NewEmptyStatus returns a new empty (blank) status
func (value *NamedValue) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &Service_NamedValue_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (value *NamedValue) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(value.Spec)
	return value.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (value *NamedValue) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*Service_NamedValue_STATUS); ok {
		value.Status = *st
		return nil
	}

	// Convert status to required version
	var st Service_NamedValue_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	value.Status = st
	return nil
}

// Hub marks that this NamedValue is the hub type for conversion
func (value *NamedValue) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (value *NamedValue) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: value.Spec.OriginalVersion,
		Kind:    "NamedValue",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20220801.NamedValue
// Generator information:
// - Generated from: /apimanagement/resource-manager/Microsoft.ApiManagement/stable/2022-08-01/apimnamedvalues.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/namedValues/{namedValueId}
type NamedValueList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamedValue `json:"items"`
}

// Storage version of v1api20220801.Service_NamedValue_Spec
type Service_NamedValue_Spec struct {
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern="^[^*#&+:<>?]+$"
	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName       string                            `json:"azureName,omitempty"`
	DisplayName     *string                           `json:"displayName,omitempty"`
	KeyVault        *KeyVaultContractCreateProperties `json:"keyVault,omitempty"`
	OriginalVersion string                            `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a apimanagement.azure.com/Service resource
	Owner       *genruntime.KnownResourceReference `group:"apimanagement.azure.com" json:"owner,omitempty" kind:"Service"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	Secret      *bool                              `json:"secret,omitempty"`
	Tags        []string                           `json:"tags,omitempty"`
	Value       *string                            `json:"value,omitempty"`
}

var _ genruntime.ConvertibleSpec = &Service_NamedValue_Spec{}

// ConvertSpecFrom populates our Service_NamedValue_Spec from the provided source
func (value *Service_NamedValue_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == value {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(value)
}

// ConvertSpecTo populates the provided destination from our Service_NamedValue_Spec
func (value *Service_NamedValue_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == value {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(value)
}

// Storage version of v1api20220801.Service_NamedValue_STATUS
type Service_NamedValue_STATUS struct {
	Conditions  []conditions.Condition             `json:"conditions,omitempty"`
	DisplayName *string                            `json:"displayName,omitempty"`
	Id          *string                            `json:"id,omitempty"`
	KeyVault    *KeyVaultContractProperties_STATUS `json:"keyVault,omitempty"`
	Name        *string                            `json:"name,omitempty"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	Secret      *bool                              `json:"secret,omitempty"`
	Tags        []string                           `json:"tags,omitempty"`
	Type        *string                            `json:"type,omitempty"`
	Value       *string                            `json:"value,omitempty"`
}

var _ genruntime.ConvertibleStatus = &Service_NamedValue_STATUS{}

// ConvertStatusFrom populates our Service_NamedValue_STATUS from the provided source
func (value *Service_NamedValue_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == value {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(value)
}

// ConvertStatusTo populates the provided destination from our Service_NamedValue_STATUS
func (value *Service_NamedValue_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == value {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(value)
}

// Storage version of v1api20220801.KeyVaultContractCreateProperties
// Create keyVault contract details.
type KeyVaultContractCreateProperties struct {
	IdentityClientId           *string                        `json:"identityClientId,omitempty" optionalConfigMapPair:"IdentityClientId"`
	IdentityClientIdFromConfig *genruntime.ConfigMapReference `json:"identityClientIdFromConfig,omitempty" optionalConfigMapPair:"IdentityClientId"`
	PropertyBag                genruntime.PropertyBag         `json:"$propertyBag,omitempty"`
	SecretIdentifier           *string                        `json:"secretIdentifier,omitempty"`
}

// Storage version of v1api20220801.KeyVaultContractProperties_STATUS
// KeyVault contract details.
type KeyVaultContractProperties_STATUS struct {
	IdentityClientId *string                                            `json:"identityClientId,omitempty"`
	LastStatus       *KeyVaultLastAccessStatusContractProperties_STATUS `json:"lastStatus,omitempty"`
	PropertyBag      genruntime.PropertyBag                             `json:"$propertyBag,omitempty"`
	SecretIdentifier *string                                            `json:"secretIdentifier,omitempty"`
}

// Storage version of v1api20220801.KeyVaultLastAccessStatusContractProperties_STATUS
// Issue contract Update Properties.
type KeyVaultLastAccessStatusContractProperties_STATUS struct {
	Code         *string                `json:"code,omitempty"`
	Message      *string                `json:"message,omitempty"`
	PropertyBag  genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	TimeStampUtc *string                `json:"timeStampUtc,omitempty"`
}

func init() {
	SchemeBuilder.Register(&NamedValue{}, &NamedValueList{})
}
