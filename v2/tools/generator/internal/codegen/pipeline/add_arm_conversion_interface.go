/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package pipeline

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/tools/generator/internal/armconversion"
	"github.com/Azure/azure-service-operator/v2/tools/generator/internal/astmodel"
	"github.com/Azure/azure-service-operator/v2/tools/generator/internal/config"
)

// ApplyARMConversionInterfaceStageID is the unique identifier of this pipeline stage
const ApplyARMConversionInterfaceStageID = "applyArmConversionInterface"

// ApplyARMConversionInterface adds the genruntime.ARMTransformer interface and the Owner property
// to all Kubernetes types.
// The genruntime.ARMTransformer interface is used to convert from the Kubernetes type to the corresponding ARM type and back.
func ApplyARMConversionInterface(idFactory astmodel.IdentifierFactory, config *config.ObjectModelConfiguration) *Stage {
	return NewLegacyStage(
		ApplyARMConversionInterfaceStageID,
		"Add ARM conversion interfaces to Kubernetes types",
		func(ctx context.Context, definitions astmodel.TypeDefinitionSet) (astmodel.TypeDefinitionSet, error) {
			converter := &armConversionApplier{
				definitions: definitions,
				idFactory:   idFactory,
				config:      config,
			}

			return converter.transformTypes()
		})
}

// GetARMTypeDefinition gets the ARM type definition for a given Kubernetes type name.
// If no matching definition can be found an error is returned.
func GetARMTypeDefinition(defs astmodel.TypeDefinitionSet, name astmodel.InternalTypeName) (astmodel.TypeDefinition, error) {
	armDefinition, ok := defs[astmodel.CreateARMTypeName(name)]
	if !ok {
		return astmodel.TypeDefinition{}, errors.Errorf("couldn't find ARM definition matching kube name %q", name)
	}

	return armDefinition, nil
}

type armConversionApplier struct {
	definitions astmodel.TypeDefinitionSet
	idFactory   astmodel.IdentifierFactory
	config      *config.ObjectModelConfiguration
}

// transformResourceSpecs applies the genruntime.ARMTransformer interface to all resource Spec types.
// It also adds the Owner property.
func (c *armConversionApplier) transformResourceSpecs() (astmodel.TypeDefinitionSet, error) {
	result := make(astmodel.TypeDefinitionSet)

	resources := c.definitions.Where(func(td astmodel.TypeDefinition) bool {
		_, ok := astmodel.AsResourceType(td.Type())
		return ok
	})

	for _, td := range resources {
		resource, ok := astmodel.AsResourceType(td.Type())
		if !ok {
			return nil, errors.Errorf("%q was not a resource, instead %T", td.Name(), td.Type())
		}

		specDefinition, err := c.transformSpec(resource)
		if err != nil {
			return nil, err
		}

		armSpecDefinition, err := GetARMTypeDefinition(c.definitions, specDefinition.Name())
		if err != nil {
			return nil, err
		}

		specDefinition, err = c.addARMConversionInterface(specDefinition, armSpecDefinition, armconversion.TypeKindSpec)
		if err != nil {
			return nil, err
		}

		result.Add(specDefinition)
	}

	return result, nil
}

// transformResourceStatuses applies the genruntime.ARMTransformer interface to all resource Status types.
func (c *armConversionApplier) transformResourceStatuses() (astmodel.TypeDefinitionSet, error) {
	result := make(astmodel.TypeDefinitionSet)

	statusDefs := c.definitions.Where(func(def astmodel.TypeDefinition) bool {
		_, ok := astmodel.AsObjectType(def.Type())
		// TODO: We need labels
		// Some status types are initially anonymous and then get named later (so end with a _Status_Xyz suffix)
		return ok && def.Name().IsStatus() && !astmodel.ARMFlag.IsOn(def.Type())
	})

	for _, td := range statusDefs {
		statusType := astmodel.IgnoringErrors(td.Type())
		if statusType == nil {
			continue
		}

		armStatusDefinition, err := GetARMTypeDefinition(c.definitions, td.Name())
		if err != nil {
			return nil, err
		}

		statusDefinition, err := c.addARMConversionInterface(td, armStatusDefinition, armconversion.TypeKindStatus)
		if err != nil {
			return nil, err
		}

		result.Add(statusDefinition)
	}

	return result, nil
}

// transformTypes adds the required ARM conversion information to all applicable types.
// If a type doesn't need any modification, it is returned unmodified.
func (c *armConversionApplier) transformTypes() (astmodel.TypeDefinitionSet, error) {
	result := make(astmodel.TypeDefinitionSet)

	// Specs
	specs, err := c.transformResourceSpecs()
	if err != nil {
		return nil, err
	}
	result.AddTypes(specs)

	// Status
	statuses, err := c.transformResourceStatuses()
	if err != nil {
		return nil, err
	}
	result.AddTypes(statuses)

	// Everything else
	otherDefs := c.definitions.Except(result)
	for _, td := range otherDefs {

		_, isObjectType := astmodel.AsObjectType(td.Type())
		hasARMFlag := astmodel.ARMFlag.IsOn(td.Type())
		requiresARMType := requiresARMType(td)
		if !isObjectType || hasARMFlag || !requiresARMType {
			// No special handling needed just add the existing type and continue
			result.Add(td)
			continue
		}

		armDefinition, err := GetARMTypeDefinition(c.definitions, td.Name())
		if err != nil {
			return nil, err
		}

		modifiedDef, err := c.addARMConversionInterface(td, armDefinition, armconversion.TypeKindOrdinary)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to add ARM conversion interface to %q", td.Name())
		}

		result.Add(modifiedDef)
	}

	return result, nil
}

// transformSpec adds an owner property to the given resource spec, and adds the ARM conversion interface
// to the spec with some special property remappings (for Type, Name, APIVersion, etc).
func (c *armConversionApplier) transformSpec(resourceType *astmodel.ResourceType) (astmodel.TypeDefinition, error) {
	resourceSpecDef, err := c.definitions.ResolveResourceSpecDefinition(resourceType)
	if err != nil {
		return astmodel.TypeDefinition{}, err
	}

	injectOwnerProperty := func(t *astmodel.ObjectType) (*astmodel.ObjectType, error) {
		if !resourceType.Owner().IsEmpty() && resourceType.Scope() == astmodel.ResourceScopeResourceGroup {
			ownerProperty := c.createOwnerProperty(resourceType.Owner())
			t = t.WithProperty(ownerProperty)
		} else if resourceType.Scope() == astmodel.ResourceScopeExtension {
			t = t.WithProperty(c.createExtensionResourceOwnerProperty())
		}

		return t, nil
	}

	remapProperties := func(t *astmodel.ObjectType) (*astmodel.ObjectType, error) {
		// TODO: Right now the Kubernetes type has all of its standard requiredness (validations). If we want to allow
		// TODO: users to submit "just a name and owner" types we will have to strip some validation until
		// TODO: https://github.com/kubernetes-sigs/controller-tools/issues/461 is fixed

		nameProp, hasName := t.Property(astmodel.NameProperty)
		if !hasName {
			return t, nil
		}

		// rename Name to AzureName and promote type if needed
		// Note: if this type ends up wrapped in another type we may need to use a visitor to do this instead of
		// doing it manually.
		namePropType := nameProp.PropertyType()
		if optional, ok := namePropType.(*astmodel.OptionalType); ok {
			namePropType = optional.Element()
		}
		azureNameProp := armconversion.GetAzureNameProperty(c.idFactory).WithType(namePropType)
		return t.WithoutProperty(astmodel.NameProperty).WithProperty(azureNameProp), nil
	}

	kubernetesDef, err := resourceSpecDef.ApplyObjectTransformations(remapProperties, injectOwnerProperty)
	if err != nil {
		return astmodel.TypeDefinition{}, errors.Wrapf(err, "remapping properties of Kubernetes definition")
	}

	return kubernetesDef, nil
}

func (c *armConversionApplier) addARMConversionInterface(
	kubeDef astmodel.TypeDefinition,
	armDef astmodel.TypeDefinition,
	typeType armconversion.TypeKind,
) (astmodel.TypeDefinition, error) {
	objectType, ok := astmodel.AsObjectType(armDef.Type())
	emptyDef := astmodel.TypeDefinition{}
	if !ok {
		return emptyDef, errors.Errorf("ARM definition %q did not define an object type", armDef.Name())
	}

	// Determine if we need special handling for collection properties. Some RPs we need to send empty collections rather
	// than nil collections, we need to determine if this def is subject to this requirement
	payloadType, err := c.config.PayloadType.Lookup(kubeDef.Name().InternalPackageReference())
	if err != nil {
		if config.IsNotConfiguredError(err) {
			// Default to 'omitempty' if not configured
			payloadType = config.OmitEmptyProperties
		} else {
			// otherwise we return an error
			return emptyDef, errors.Wrapf(err, "looking up payload type for %q", kubeDef.Name())
		}
	}

	addInterfaceHandler := func(t *astmodel.ObjectType) (astmodel.Type, error) {
		result := t.WithInterface(armconversion.NewARMConversionImplementation(
			armDef.Name(),
			objectType,
			c.idFactory,
			typeType,
			payloadType))
		return result, nil
	}

	result, err := kubeDef.ApplyObjectTransformation(addInterfaceHandler)
	if err != nil {
		emptyDef := astmodel.TypeDefinition{}
		return emptyDef,
			errors.Errorf("failed to add ARM conversion interface to Kubenetes object definition %s", armDef.Name())
	}

	return result, nil
}

func (c *armConversionApplier) createOwnerProperty(ownerTypeName astmodel.InternalTypeName) *astmodel.PropertyDefinition {
	grp := ownerTypeName.InternalPackageReference().Group()
	group := grp + astmodel.GroupSuffix
	kind := ownerTypeName.Name()

	prop := astmodel.NewPropertyDefinition(
		c.idFactory.CreatePropertyName(astmodel.OwnerProperty, astmodel.Exported),
		c.idFactory.CreateStringIdentifier(astmodel.OwnerProperty, astmodel.NotExported),
		astmodel.OptionalKnownResourceReferenceType)
	prop = prop.WithDescription(
		fmt.Sprintf("The owner of the resource. The owner controls where the resource goes when it is deployed. "+
			"The owner also controls the resources lifecycle. "+
			"When the owner is deleted the resource will also be deleted. Owner is expected to "+
			"be a reference to a %s/%s resource", group, kind))
	prop = prop.WithTag("group", group)
	prop = prop.WithTag("kind", kind)
	prop = prop.MakeRequired() // Owner is always required

	return prop
}

func (c *armConversionApplier) createExtensionResourceOwnerProperty() *astmodel.PropertyDefinition {
	prop := astmodel.NewPropertyDefinition(
		c.idFactory.CreatePropertyName(astmodel.OwnerProperty, astmodel.Exported),
		c.idFactory.CreateStringIdentifier(astmodel.OwnerProperty, astmodel.NotExported),
		astmodel.NewOptionalType(astmodel.ArbitraryOwnerReference)).MakeRequired()
	prop = prop.WithDescription(
		"The owner of the resource. The owner controls where the resource goes when it is deployed. " +
			"The owner also controls the resources lifecycle. " +
			"When the owner is deleted the resource will also be deleted. " +
			"This resource is an extension resource, which means that any other Azure resource can be its owner.")
	return prop
}
