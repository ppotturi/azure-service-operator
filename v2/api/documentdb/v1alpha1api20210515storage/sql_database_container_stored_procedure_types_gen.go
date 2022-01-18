// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package v1alpha1api20210515storage

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +kubebuilder:rbac:groups=documentdb.azure.com,resources=sqldatabasecontainerstoredprocedures,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=documentdb.azure.com,resources={sqldatabasecontainerstoredprocedures/status,sqldatabasecontainerstoredprocedures/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
//Storage version of v1alpha1api20210515.SqlDatabaseContainerStoredProcedure
//Generated from: https://schema.management.azure.com/schemas/2021-05-15/Microsoft.DocumentDB.json#/resourceDefinitions/databaseAccounts_sqlDatabases_containers_storedProcedures
type SqlDatabaseContainerStoredProcedure struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec `json:"spec,omitempty"`
	Status            SqlStoredProcedureGetResults_Status                         `json:"status,omitempty"`
}

var _ conditions.Conditioner = &SqlDatabaseContainerStoredProcedure{}

// GetConditions returns the conditions of the resource
func (procedure *SqlDatabaseContainerStoredProcedure) GetConditions() conditions.Conditions {
	return procedure.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (procedure *SqlDatabaseContainerStoredProcedure) SetConditions(conditions conditions.Conditions) {
	procedure.Status.Conditions = conditions
}

var _ genruntime.KubernetesResource = &SqlDatabaseContainerStoredProcedure{}

// AzureName returns the Azure name of the resource
func (procedure *SqlDatabaseContainerStoredProcedure) AzureName() string {
	return procedure.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2021-05-15"
func (procedure SqlDatabaseContainerStoredProcedure) GetAPIVersion() string {
	return "2021-05-15"
}

// GetResourceKind returns the kind of the resource
func (procedure *SqlDatabaseContainerStoredProcedure) GetResourceKind() genruntime.ResourceKind {
	return genruntime.ResourceKindNormal
}

// GetSpec returns the specification of this resource
func (procedure *SqlDatabaseContainerStoredProcedure) GetSpec() genruntime.ConvertibleSpec {
	return &procedure.Spec
}

// GetStatus returns the status of this resource
func (procedure *SqlDatabaseContainerStoredProcedure) GetStatus() genruntime.ConvertibleStatus {
	return &procedure.Status
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/storedProcedures"
func (procedure *SqlDatabaseContainerStoredProcedure) GetType() string {
	return "Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers/storedProcedures"
}

// NewEmptyStatus returns a new empty (blank) status
func (procedure *SqlDatabaseContainerStoredProcedure) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &SqlStoredProcedureGetResults_Status{}
}

// Owner returns the ResourceReference of the owner, or nil if there is no owner
func (procedure *SqlDatabaseContainerStoredProcedure) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(procedure.Spec)
	return &genruntime.ResourceReference{
		Group: group,
		Kind:  kind,
		Name:  procedure.Spec.Owner.Name,
	}
}

// SetStatus sets the status of this resource
func (procedure *SqlDatabaseContainerStoredProcedure) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*SqlStoredProcedureGetResults_Status); ok {
		procedure.Status = *st
		return nil
	}

	// Convert status to required version
	var st SqlStoredProcedureGetResults_Status
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	procedure.Status = st
	return nil
}

// Hub marks that this SqlDatabaseContainerStoredProcedure is the hub type for conversion
func (procedure *SqlDatabaseContainerStoredProcedure) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (procedure *SqlDatabaseContainerStoredProcedure) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: procedure.Spec.OriginalVersion,
		Kind:    "SqlDatabaseContainerStoredProcedure",
	}
}

// +kubebuilder:object:root=true
//Storage version of v1alpha1api20210515.SqlDatabaseContainerStoredProcedure
//Generated from: https://schema.management.azure.com/schemas/2021-05-15/Microsoft.DocumentDB.json#/resourceDefinitions/databaseAccounts_sqlDatabases_containers_storedProcedures
type SqlDatabaseContainerStoredProcedureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SqlDatabaseContainerStoredProcedure `json:"items"`
}

//Storage version of v1alpha1api20210515.DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec
type DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec struct {
	//AzureName: The name of the resource in Azure. This is often the same as the name
	//of the resource in Kubernetes but it doesn't have to be.
	AzureName       string               `json:"azureName"`
	Location        *string              `json:"location,omitempty"`
	Options         *CreateUpdateOptions `json:"options,omitempty"`
	OriginalVersion string               `json:"originalVersion"`

	// +kubebuilder:validation:Required
	Owner       genruntime.KnownResourceReference `group:"documentdb.azure.com" json:"owner" kind:"SqlDatabaseContainer"`
	PropertyBag genruntime.PropertyBag            `json:"$propertyBag,omitempty"`
	Resource    *SqlStoredProcedureResource       `json:"resource,omitempty"`
	Tags        map[string]string                 `json:"tags,omitempty"`
}

var _ genruntime.ConvertibleSpec = &DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec{}

// ConvertSpecFrom populates our DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec from the provided source
func (procedures *DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == procedures {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(procedures)
}

// ConvertSpecTo populates the provided destination from our DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec
func (procedures *DatabaseAccountsSqlDatabasesContainersStoredProcedures_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == procedures {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(procedures)
}

//Storage version of v1alpha1api20210515.SqlStoredProcedureGetResults_Status
type SqlStoredProcedureGetResults_Status struct {
	Conditions  []conditions.Condition                           `json:"conditions,omitempty"`
	Id          *string                                          `json:"id,omitempty"`
	Location    *string                                          `json:"location,omitempty"`
	Name        *string                                          `json:"name,omitempty"`
	PropertyBag genruntime.PropertyBag                           `json:"$propertyBag,omitempty"`
	Resource    *SqlStoredProcedureGetProperties_Status_Resource `json:"resource,omitempty"`
	Tags        map[string]string                                `json:"tags,omitempty"`
	Type        *string                                          `json:"type,omitempty"`
}

var _ genruntime.ConvertibleStatus = &SqlStoredProcedureGetResults_Status{}

// ConvertStatusFrom populates our SqlStoredProcedureGetResults_Status from the provided source
func (results *SqlStoredProcedureGetResults_Status) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == results {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(results)
}

// ConvertStatusTo populates the provided destination from our SqlStoredProcedureGetResults_Status
func (results *SqlStoredProcedureGetResults_Status) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == results {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(results)
}

//Storage version of v1alpha1api20210515.SqlStoredProcedureGetProperties_Status_Resource
type SqlStoredProcedureGetProperties_Status_Resource struct {
	Body        *string                `json:"body,omitempty"`
	Etag        *string                `json:"_etag,omitempty"`
	Id          *string                `json:"id,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Rid         *string                `json:"_rid,omitempty"`
	Ts          *float64               `json:"_ts,omitempty"`
}

//Storage version of v1alpha1api20210515.SqlStoredProcedureResource
//Generated from: https://schema.management.azure.com/schemas/2021-05-15/Microsoft.DocumentDB.json#/definitions/SqlStoredProcedureResource
type SqlStoredProcedureResource struct {
	Body        *string                `json:"body,omitempty"`
	Id          *string                `json:"id,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

func init() {
	SchemeBuilder.Register(&SqlDatabaseContainerStoredProcedure{}, &SqlDatabaseContainerStoredProcedureList{})
}
