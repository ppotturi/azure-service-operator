/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package config

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	kerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/Azure/azure-service-operator/v2/internal/util/typo"
	"github.com/Azure/azure-service-operator/v2/tools/generator/internal/astmodel"
)

// TypeConfiguration contains additional information about a specific kind of resource within a version of a group and forms
// part of a hierarchy containing information to supplement the schema and swagger sources consumed by the generator.
//
// ┌──────────────────────────┐       ┌────────────────────┐       ┌──────────────────────┐       ╔═══════════════════╗       ┌───────────────────────┐
// │                          │       │                    │       │                      │       ║                   ║       │                       │
// │ ObjectModelConfiguration │───────│ GroupConfiguration │───────│ VersionConfiguration │───────║ TypeConfiguration ║───────│ PropertyConfiguration │
// │                          │1  1..n│                    │1  1..n│                      │1  1..n║                   ║1  1..n│                       │
// └──────────────────────────┘       └────────────────────┘       └──────────────────────┘       ╚═══════════════════╝       └───────────────────────┘
type TypeConfiguration struct {
	name       string
	properties map[string]*PropertyConfiguration
	advisor    *typo.Advisor
	// Configurable properties here (alphabetical, please)
	AzureGeneratedSecrets    configurable[[]string]
	DefaultAzureName         configurable[bool]
	Export                   configurable[bool]
	ExportAs                 configurable[string]
	GeneratedConfigs         configurable[map[string]string]
	Importable               configurable[bool]
	IsResource               configurable[bool]
	ManualConfigs            configurable[[]string]
	NameInNextVersion        configurable[string]
	RenameTo                 configurable[string]
	ResourceEmbeddedInParent configurable[string]
	SupportedFrom            configurable[string]
}

const (
	azureGeneratedSecretsTag    = "$azureGeneratedSecrets"    // A set of strings specifying which secrets are generated by Azure
	generatedConfigsTag         = "$generatedConfigs"         // A map of strings specifying which spec or status properties should be exported to configmap
	manualConfigsTag            = "$manualConfigs"            // A set of strings specifying which config map fields should be generated (to be filled out by resource extension)
	exportTag                   = "$export"                   // Boolean specifying whether a resource type is exported
	exportAsTag                 = "$exportAs"                 // String specifying the name to use for a type (implies $export: true)
	importableTag               = "$importable"               // Boolean specifying whether a resource type is importable via asoctl (defaults to true)
	isResourceTag               = "$isResource"               // Boolean specifying whether a particular type is a resource or not.
	nameInNextVersionTag        = "$nameInNextVersion"        // String specifying a type or property name change in the next version
	supportedFromTag            = "$supportedFrom"            // Label specifying the first ASO release supporting the resource
	renameTo                    = "$renameTo"                 // String specifying the new name of a type
	resourceEmbeddedInParentTag = "$resourceEmbeddedInParent" // String specifying resource name of parent
	defaultAzureNameTag         = "$defaultAzureName"         // Boolean indicating if the resource should automatically default AzureName
)

func NewTypeConfiguration(name string) *TypeConfiguration {
	scope := "type " + name
	return &TypeConfiguration{
		name:       name,
		properties: make(map[string]*PropertyConfiguration),
		advisor:    typo.NewAdvisor(),
		// Initialize configurable properties here (alphabetical, please)
		AzureGeneratedSecrets:    makeConfigurable[[]string](azureGeneratedSecretsTag, scope),
		DefaultAzureName:         makeConfigurable[bool](defaultAzureNameTag, scope),
		Export:                   makeConfigurable[bool](exportTag, scope),
		ExportAs:                 makeConfigurable[string](exportAsTag, scope),
		Importable:               makeConfigurable[bool](importableTag, scope),
		IsResource:               makeConfigurable[bool](isResourceTag, scope),
		GeneratedConfigs:         makeConfigurable[map[string]string](generatedConfigsTag, scope),
		ManualConfigs:            makeConfigurable[[]string](manualConfigsTag, scope),
		NameInNextVersion:        makeConfigurable[string](nameInNextVersionTag, scope),
		RenameTo:                 makeConfigurable[string](renameTo, scope),
		ResourceEmbeddedInParent: makeConfigurable[string](resourceEmbeddedInParentTag, scope),
		SupportedFrom:            makeConfigurable[string](supportedFromTag, scope),
	}
}

// Add includes configuration for the specified property as a part of this type configuration
func (tc *TypeConfiguration) addProperty(name string, property *PropertyConfiguration) {
	// Indexed by lowercase name of the property to allow case-insensitive lookups
	tc.properties[strings.ToLower(name)] = property
}

// visitProperty invokes the provided visitor on the specified property if present.
// Returns a NotConfiguredError if the property is not found; otherwise whatever error is returned by the visitor.
func (tc *TypeConfiguration) visitProperty(
	property astmodel.PropertyName,
	visitor *configurationVisitor,
) error {
	pc, err := tc.findProperty(property)
	if err != nil {
		return err
	}

	err = visitor.visitProperty(pc)
	if err != nil {
		return errors.Wrapf(err, "configuration of type %s", tc.name)
	}

	return nil
}

// visitProperties invokes the provided visitor on all properties.
func (tc *TypeConfiguration) visitProperties(visitor *configurationVisitor) error {
	errs := make([]error, 0, len(tc.properties))
	for _, pc := range tc.properties {
		err := visitor.visitProperty(pc)
		err = tc.advisor.Wrapf(err, pc.name, "property %s not seen", pc.name)
		errs = append(errs, err)
	}

	// Both errors.Wrapf() and kerrors.NewAggregate() return nil if nothing went wrong
	return errors.Wrapf(
		kerrors.NewAggregate(errs),
		"type %s",
		tc.name)
}

// findProperty uses the provided property name to work out which nested PropertyConfiguration should be used
// either returns the requested property configuration, or an error saying that it couldn't be found
func (tc *TypeConfiguration) findProperty(property astmodel.PropertyName) (*PropertyConfiguration, error) {
	// Store the property id using lowercase,
	// so we can do case-insensitive lookups later
	tc.advisor.AddTerm(string(property))
	p := strings.ToLower(string(property))
	if pc, ok := tc.properties[p]; ok {
		return pc, nil
	}

	msg := fmt.Sprintf(
		"configuration of type %s has no detail for property %s",
		tc.name,
		property)
	return nil, NewNotConfiguredError(msg).WithOptions("properties", tc.configuredProperties())
}

// UnmarshalYAML populates our instance from the YAML.
// The slice node.Content contains pairs of nodes, first one for an ID, then one for the value.
func (tc *TypeConfiguration) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return errors.New("expected mapping")
	}

	tc.properties = make(map[string]*PropertyConfiguration)
	var lastId string

	for i, c := range value.Content {
		// Grab identifiers and loop to handle the associated value
		if i%2 == 0 {
			lastId = c.Value
			continue
		}

		if strings.EqualFold(lastId, generatedConfigsTag) && c.Kind == yaml.MappingNode {
			azureGeneratedConfigs := make(map[string]string)

			idx := 0
			var key string
			for _, content := range c.Content {
				var val string
				if content.Kind == yaml.ScalarNode {
					val = content.Value
				} else {
					return errors.Errorf(
						"unexpected yam value for %s (line %d col %d)",
						generatedConfigsTag,
						content.Line,
						content.Column)
				}

				// first value is the key
				if idx%2 == 0 {
					key = val
				}
				if idx%2 == 1 {
					if !strings.HasPrefix(val, "$.") {
						return errors.Errorf("%s entry %q must begin with $.", generatedConfigsTag, val)
					}
					azureGeneratedConfigs[key] = val

				}
				idx++
			}

			// TODO: Check we had an even number of nodes

			tc.GeneratedConfigs.Set(azureGeneratedConfigs)
			continue
		}

		// Handle nested property metadata
		if c.Kind == yaml.MappingNode {
			p := NewPropertyConfiguration(lastId)
			err := c.Decode(p)
			if err != nil {
				return errors.Wrapf(err, "decoding yaml for %q", lastId)
			}

			tc.addProperty(lastId, p)
			continue
		}

		// $nameInNextVersion: <string>
		if strings.EqualFold(lastId, nameInNextVersionTag) && c.Kind == yaml.ScalarNode {
			tc.NameInNextVersion.Set(c.Value)
			continue
		}

		// $export: <bool>
		if strings.EqualFold(lastId, exportTag) && c.Kind == yaml.ScalarNode {
			var export bool
			err := c.Decode(&export)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", exportTag)
			}

			tc.Export.Set(export)
			continue
		}

		// $exportAs: <string>
		if strings.EqualFold(lastId, exportAsTag) && c.Kind == yaml.ScalarNode {
			tc.ExportAs.Set(c.Value)
			continue
		}

		// $azureGeneratedSecrets:
		// - secret1
		// - secret2
		if strings.EqualFold(lastId, azureGeneratedSecretsTag) && c.Kind == yaml.SequenceNode {
			var azureGeneratedSecrets []string
			for _, content := range c.Content {
				if content.Kind == yaml.ScalarNode {
					azureGeneratedSecrets = append(azureGeneratedSecrets, content.Value)
				} else {
					return errors.Errorf(
						"unexpected yam value for %s (line %d col %d)",
						azureGeneratedSecretsTag,
						content.Line,
						content.Column)
				}
			}

			tc.AzureGeneratedSecrets.Set(azureGeneratedSecrets)
			continue
		}

		// $manualConfigs
		// - config1
		// - config2
		if strings.EqualFold(lastId, manualConfigsTag) && c.Kind == yaml.SequenceNode {
			var manualAzureGeneratedConfigs []string
			for _, content := range c.Content {
				if content.Kind == yaml.ScalarNode {
					manualAzureGeneratedConfigs = append(manualAzureGeneratedConfigs, content.Value)
				} else {
					return errors.Errorf(
						"unexpected yam value for %s (line %d col %d)",
						manualConfigsTag,
						content.Line,
						content.Column)
				}
			}

			tc.ManualConfigs.Set(manualAzureGeneratedConfigs)
			continue
		}

		// $SupportedFrom
		if strings.EqualFold(lastId, supportedFromTag) && c.Kind == yaml.ScalarNode {
			tc.SupportedFrom.Set(c.Value)
			continue
		}

		// $renameTo: <string>
		if strings.EqualFold(lastId, renameTo) && c.Kind == yaml.ScalarNode {
			var renameTo string
			err := c.Decode(&renameTo)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", renameTo)
			}

			tc.RenameTo.Set(renameTo)
			continue
		}

		// $resourceEmbeddedInParent: <string>
		if strings.EqualFold(lastId, resourceEmbeddedInParentTag) && c.Kind == yaml.ScalarNode {
			var resourceEmbeddedInParent string
			err := c.Decode(&resourceEmbeddedInParent)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", resourceEmbeddedInParentTag)
			}

			tc.ResourceEmbeddedInParent.Set(resourceEmbeddedInParent)
			continue
		}

		// $isResource: <bool>
		if strings.EqualFold(lastId, isResourceTag) && c.Kind == yaml.ScalarNode {
			var isResource bool
			err := c.Decode(&isResource)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", isResourceTag)
			}

			tc.IsResource.Set(isResource)
			continue
		}

		// $importable: <bool>
		if strings.EqualFold(lastId, importableTag) && c.Kind == yaml.ScalarNode {
			var importable bool
			err := c.Decode(&importable)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", importableTag)
			}

			tc.Importable.Set(importable)
			continue
		}

		// $defaultAzureName: <bool>
		if strings.EqualFold(lastId, defaultAzureNameTag) && c.Kind == yaml.ScalarNode {
			var defaultAzureName bool
			err := c.Decode(&defaultAzureName)
			if err != nil {
				return errors.Wrapf(err, "decoding %s", defaultAzureNameTag)
			}

			tc.DefaultAzureName.Set(defaultAzureName)
			continue
		}

		// No handler for this value, return an error
		return errors.Errorf(
			"type configuration, unexpected yaml value %s: %s (line %d col %d)", lastId, c.Value, c.Line, c.Column)
	}

	return nil
}

// configuredProperties returns a sorted slice containing all the properties configured on this type
func (tc *TypeConfiguration) configuredProperties() []string {
	result := make([]string, 0, len(tc.properties))
	for _, c := range tc.properties {
		// Use the actual names of the properties, not the lower-cased keys of the map
		result = append(result, c.name)
	}

	return result
}
