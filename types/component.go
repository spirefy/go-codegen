package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type (
	Components      []*Component
	ComponentSource int32
)

// Types of Source locations where a Component is found/defined/referenced
const (
	SourceUnknown            ComponentSource = iota // Anything else or possibly default initialized value
	SourceComponent                                 // Source is from a defined component.. e.g. one now found in a parameter, request, response, or is a reference, inlined, or part of a property.
	SourceParameter                                 // Request parameters that define a component
	SourceInline                                    // Inline component.. e.g. defined in place rather than a reference to a defined component
	SourceReference                                 // Reference to a component
	SourceRequestBodyInline                         // Request body inlined component
	SourceResponseBodyInline                        // Response body inlined component
	SourceProperty                                  // Property of a component.. should be used for properties of a component
)

// This is a generic component structure.. tries to capture all possible pieces of data any sort of component might contain.. a superset of different component implementations if you will.
type Component struct {
	// sync.Mutex
	Id          int             // This is used to allow unique ids per component, mostly to be used to reference and used to find from references.
	Name        string          // name of this component
	RawName     string          // this is the name as it appears in the source document. Name refers to the name as it should be used when generating output and could be modified using extensions
	Type        string          // matches with json schema 2020-12 types.. object, array, string, number, enum. When number, Format will contain the sub type (int, int32, float64, etc)
	Description string          // A description of this component if available
	Format      string          // if this is a direct component.. not object.. specifies the format of the type if provided (e.g. int64, in32, float64, phone, email, etc.. )
	Required    *bool           // true if this component is required (likely used for validation purposes at runtime)
	Null        *bool           // true if this component can be null, false if this can not be null
	Enums       []string        // A property may be an enum type
	Source      ComponentSource // Can be used by loaders to indicate the source of this component.. is it part of a request or response inline body, defined component, other? (use 'inline' for an inline component, 'component' if defined)
	Raw         json.RawMessage // The raw Json of this particular property
	Ref         any             // This is "any" so that if the type is Object this would either be a Component ref or a string name placeholder (until can be resolved after all components are processed by all loaders). If type is an array, this is a string that holds the primitive type of array
	SourceDoc   string          // This refs the URL/path (or alias/name) to the doc that this component came from. This is particularly usefule when trying to merge two (or more) similar components in to one.. to ensure they are from the same doc.. as it is possible for two (or more) different APIs from different organizations to be loaded
	Version     string          // This is the version of this component. Possible multiple versions of same component might be loaded via multiple sources.
	Latest      bool            // Indicates that this is the latest version of a component based on the version value
	Properties  Properties      // if this component has any associated properties, this contains the slice of those properties
}

// Len
// Part of the sorting interface implementation for custom sorting Components
func (c Components) Len() int {
	return len(c)
}

// Swap
// Part of the sorting interface implementation for custom sorting Components
func (c Components) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less
// Part of the sorting interface implementation for custom sorting Components
func (c Components) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

func (cs ComponentSource) String() string {
	switch cs {
	case SourceComponent:
		return "Defined Component"
	case SourceInline:
		return "Inlined Component"
	case SourceParameter:
		return "Parameter Component"
	case SourceProperty:
		return "Property Component"
	case SourceReference:
		return "Reference Component"
	case SourceRequestBodyInline:
		return "Request Body Inlined Component"
	case SourceResponseBodyInline:
		return "Response Body Inlined Component"
	case SourceUnknown:
		return "Unknown Component Type"
	default:
		return ""
	}
}

func (c Components) FindComponentById(id int) *Component {
	for _, comp := range c {
		if comp.Id == id {
			return comp
		}
	}

	return nil
}

func (c Components) FindComponentByName(name string) *Component {
	if nil != c {
		for _, comp := range c {
			if strings.EqualFold(comp.Name, name) {
				return comp
			}
		}
	}

	return nil
}

func (c Components) FindComponentsByIds(components Components) []*Component {
	found := make([]*Component, 0)

	for _, comp := range components {
		component := c.FindComponentById(comp.Id)

		if nil != component {
			found = append(found, component)
		}
	}

	return found
}

// FindComponentByComparison
//
// This method will iterate the range of receiver Components to see if it can find
// the provided component. It does so by comparing several fields. If a match is found
// a pointer to the component is returned otherwise nil is returned
func (c Components) FindComponentByComparison(component *Component) *Component {
	if nil != c && nil != component {
		for _, comp := range c {
			if strings.EqualFold(comp.Name, component.Name) &&
				comp.Source == component.Source &&
				comp.Format == component.Format &&
				comp.Required == component.Required &&
				// This check is to ensure that the sourceDoc name/alias/path is the SAME as that of the provided Component, ensuring they are from the same source API
				(len(comp.SourceDoc) > 0 && len(component.SourceDoc) > 0 && strings.ToLower(comp.SourceDoc) == strings.ToLower(component.SourceDoc)) {
				return comp
			}
		}
	}

	return nil
}

// GetDefinedComponents
//
// This method will return any Component objects that have a type of SourceComponent as Source.. which indicates
// a defined component. A Defined component is one where it is outside the scope of an inlined component found in
// resource request/response bodies or parameters. In other words.. it will most certainly have a name to it and is
// referenced elsewhere in API definitions so as to be reusable. Inlined components are "dynamic" in that while they
// may be duplicate definitions.. they are often custom for the specific location they are used inline.
//
// TODO: Replace this with a static slice of components to avoid having to create/loop every call to this function. It
// should be created during the parsing stages of source loading.
func (c Components) GetDefinedComponents() Components {
	components := make(Components, 0)

	for _, comp := range c {
		if comp.Source == SourceComponent {
			components = append(components, comp)
		}
	}

	return components
}

// GetParameterComponents
//
// This method will return any components defined as parameters.. namely a request parameter
// What is a parameter component? Basically it's a component defined as part of a request parameter typically as a json string value.
// The source type is set to SourceParameter by a loader if it deems an object rather than a primitive type is a parameter value.
func (c Components) GetParameterComponents() Components {
	components := make(Components, 0)

	for _, comp := range c {
		if comp.Source == SourceParameter {
			components = append(components, comp)
		}
	}

	return components
}

// GetInlinedComponents
// This method will iterate all components and return only those that are
// found inline.. e.g. defined in place rather than a reference to a "defined" component.
// These will typically be request and response payloads that are possibly one off structures or
// request parameters that accept a json object
func (c Components) GetInlinedComponents() Components {
	components := make(Components, 0)

	for _, comp := range c {
		if comp.Source == SourceInline || comp.Source == SourceRequestBodyInline || comp.Source == SourceResponseBodyInline {
			components = append(components, comp)
		}
	}

	return components
}

// NewComponent
//
// This method will attempt to create a new component from the provided parameters. The purpose of having all these parameters rather than a single Component object is to ensure any
// required fields, or other conditions (e.g. if version is added make sure same version doesn't already exist.. etc) before adding the new Component to the *Components receiver slice.
//
// It will attempt to find a matching component by comparing name, version, source, type and format. If a match is found, the matched object is returned and the calling function or method
// should utilize the components Id field for referencing. The Id field is auto incremented if it's a new component.
func (c *Components) NewComponent(name, rawName, typ, description, format, version string, required, null *bool, latest bool, enums []string, source ComponentSource, ref any, raw json.RawMessage, replaceOrMerge bool) (component *Component, err error) {
	// make sure some must have fields are not nil/emtpy
	if len(name) <= 0 {
		err = fmt.Errorf(" component must have a name ")
	} else {
		component = &Component{
			Id:          generateUniqueInt(),
			Name:        name,
			RawName:     rawName,
			Type:        typ,
			Description: description,
			Format:      format,
			Source:      source,
			Version:     version,
			Required:    required,
			Null:        null,
			Latest:      latest,
			Enums:       enums,
			Ref:         ref,
			Raw:         raw,
		}

		// Look up to see if comp exists
		cmp := c.FindComponentByComparison(component)

		if nil != cmp {
			if replaceOrMerge {

				// comp is a NEW component.. we set the auto incrementing ID

				// TODO: MERGE component.. e.g. do we simply "ignore" the incoming component due to it matching the existing one? Or do we merge any additional
				// fields that the new incoming component (comp) may contain into the already existing component?

				// Add stat for duplicate (merged) component
				// TODO: Log the original component and new component source details
				// TODO: for more stats infos
			}
		} else {
			*c = append(*c, component)
		}
	}

	// Sort all the components
	sort.Sort(c)
	return
}
