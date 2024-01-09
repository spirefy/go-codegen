package types

import (
	"encoding/json"
	"fmt"
)

type (
	Properties []*Property
)

type Property struct {
	Id          int             // This is used to allow unique ids per component, mostly to be used to reference and used to find from references.
	Name        string          // name of this component
	RawName     string          // this is the name as it appears in the source. Name would indicate the raw name or could be a generated name as per the loader processing it. This property should be as it appears in the source.
	Type        string          // matches with json schema 2020-12 types.. object, array, string, number, enum. When number, Format will contain the sub type (int, int32, float64, etc)
	Description string          // A description of this component if available
	Format      string          // if this is a direct component.. not object.. specifies the format of the type if provided (e.g. int64, in32, float64, phone, email, etc.. )
	Required    *bool           // true if this component is required (likely used for validation purposes at runtime)
	Null        *bool           // true if this component can be null, false if this can not be null
	Enums       []string        // A property may be an enum type
	Raw         json.RawMessage // The raw Json of this particular property
	Properties  Properties      // if this property has any associated properties, this contains the slice of those properties
	Version     string          // This is the version of this component. Possible multiple versions of same component might be loaded via multiple sources.
	Latest      bool            // Indicates that this is the latest version of a component based on the version value
	Ref         any             // This is "any" so that if the type is Object this would either be a Component ref or a string name placeholder (until can be resolved after all components are processed by all loaders). If type is an array, this is a string that holds the primitive type of array
}

// Len
// Part of the sorting interface implementation for custom sorting Properties
func (p Properties) Len() int {
	return len(p)
}

// Swap
// Part of the sorting interface implementation for custom sorting Properties
func (p Properties) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less
// Part of the sorting interface implementation for custom sorting Properties
func (p Properties) Less(i, j int) bool {
	return p[i].Name < p[j].Name
}

// NewProperty
//
// This method will attempt to create a new component from the provided parameters. The purpose of having all these parameters rather than a single Component object is to ensure any
// required fields, or other conditions (e.g. if version is added make sure same version doesn't already exist.. etc) before adding the new Component to the *Components receiver slice.
//
// It will attempt to find a matching component by comparing name, version, source, type and format. If a match is found, the matched object is returned and the calling function or method
// should utilize the components Id field for referencing. The Id field is auto incremented if it's a new component.
func NewProperty(name, rawName, typ, description, format, version string, required, null *bool, latest bool, enums []string, ref any, raw json.RawMessage) (*Property, error) {
	// make sure some must have fields are not nil/empty
	if len(name) <= 0 {
		return nil, fmt.Errorf(" property must have a name ")
	}

	// if the type of property is an object, lets look up the ref

	property := &Property{
		Id:          generateUniqueInt(), // use here as well as Component to keep ids unique even across properties and components.
		Name:        name,
		Type:        typ,
		Description: description,
		Format:      format,
		Required:    required,
		Null:        null,
		Enums:       enums,
		Raw:         raw,
		Version:     version,
		Ref:         ref,
		Latest:      latest,
	}

	return property, nil
}
