package types

import (
	"fmt"
	"regexp"
	"strings"
)

type (
	Resources      []*Resource
	Parameters     []*Parameter
	Requests       []*Request
	Responses      []*Response
	ResponseBodies []*ResponseBody
)

type ResourceType int32

const (
	// HTTP REST SYNC
	HTTP      ResourceType = 1
	GRPC      ResourceType = 2
	ASYNC     ResourceType = 4
	WEBSOCKET ResourceType = 8
	GRAPHQL   ResourceType = 16
	SOAP      ResourceType = 32
	RPC       ResourceType = 64
	// FOLDER A folder means a resource has children resources and itself is a hierarchy organization marker
	FOLDER ResourceType = 1024
)

type Resource struct {
	// This is a unique id assigned at creation time. Can be used for reference if need be.
	Id int `json:"id"`

	// Path template (URL path)
	// OpenAPI Mapping - Path name
	Path string `json:"path"`

	// This is a generated UNIQUE resource ID.. it should typically be comprised of the method + the path as in an API these two combined like <method>:<path> should form a unique string value
	ResourceId string `json:"resourceId"`

	// This is the root resoure path if specified to help arrange children resources under a single root.
	//
	// For example, /users, /users/:id, /users/:id/grades would all fall under /users as the base resource path. This could
	// be used in various ways, such as when generating code.. make sure all sub resources go within a single generated file instead
	// of possibly individual files for every sub resource. Technigically the first word up to the first / after the first character
	// could be considered the Root. This is mostly a convenience IF a loader uses it.
	//
	// This should be auto set by the NewResource call.. so that every resource has this value set and generators
	// can rely on this for various uses, such as dividing code up based on the resource roots of APIs
	Root string `json:"root"`

	// HTTP method required to invoke the operation
	// OpenAPI Mapping - Path item method
	Method string `json:"method"`

	// Name for corresponding Resource. (Multiple specification has difference between name and summary also
	// webhooks are map with name as key and value as path object in OAS)
	Name string `json:"name"`

	// Detailed description about Operation.
	// OpenAPI Mapping - Operation description
	Description string `json:"description,omitempty"`

	// Short summary about Operation
	// OpenAPI Mapping - Operation summary
	Summary string `json:"summary,omitempty"`

	// Marks the operation as deprecated
	// OpenAPI Mapping - Operation deprecated
	Deprecated bool `json:"deprecated,omitempty"`

	// This can be the protocol or other use for this particular resource. For example, it can be http or grpc or websocket to represent this resource as one of those
	// types. Or.. in the case of a collection.. it can be 'folder' to represent that this is a rolder, not an actual resource, and in that case the Resources [] would
	// likely not be 0 length.
	ResourceType ResourceType `json:"resourceType,omitempty"`

	// Not sure what this is
	// TODO: FIGURE THIS OUT.. IS IT NEEDED
	Parameters Parameters `json:"parameters,omitempty"`

	// Request information required by the operation
	// OpenAPI Mapping - Operation requestBody
	Requests Requests `json:"requestBodies,omitempty"`

	// Response data returned by the operation
	// OpenAPI Mapping - Operation response
	Responses Responses `json:"responses,omitempty"`

	// An array of components this resource references
	Components *Components `json:"components,omitempty"`

	// A Map of key/value pairs, where the key is a variable name found in the corresponding source of the resource object and the
	// value is whatever the particular source loader deems as such. In the case of a Collection, a value might be found in an
	// environment associated with the collection.
	Variables map[string]string

	// This property can be set by loaders as each Resource is being created. This allows the ability for a loader to determine if
	// a resource already exists (based on method and URL being identical) and can determine if it already exists if it should
	// append/merge to it.. or overwrite it, or ignore it. For example, Postman collections can build Resource objects from
	// individual requests. A folder of requests may have duplicate method/url (resource) declared and when loading.. a collection
	// loader may first look up to see if the resource exists.. if so.. determine what the source is.. "collection" as the source
	// would mean the collection loader implementation could ignore, merge the current processing resource with the existing one
	// (e.g. merge the request bodies to build a more complete Component (type... payload.. whatever) than a single resource may
	// have access to. But if the source is say openapi or something else.. the collection loader could (should?) ignore the currently
	// processing resource as a Resource object since it already exists in a probably more complete from ... being from an API definition.
	Source string `json:"source,omitempty"`

	// This refs the URL/path (or alias/name) to the doc that this component came from. This is particularly usefule when trying to merge
	// two (or more) similar components in to one.. to ensure they are from the same doc.. as it is possible for two (or more) different APIs
	// from different organizations to be loaded
	SourceDoc string `json:"sourceDoc,omitempty"`

	// This is the version of this resource.. which should come from the source version. The purpose is to allow the lookup process
	// for a resource to identify if a resource that already exists.. is an older (or newer) version. Loader implementations that utilize
	// this field can define their own rules, but ideally it should be utilized to determine if a unique resource should be replaced.. or
	// possibly merged with the parsing operation. For example, when loading two (or more) OpenAPI definitions that have the SAME title
	// but different versions.. this value could be used to ensure the resource is of the latest version. A Postman collection loader
	// however may wish to MERGE the various properties/etc of a resource due to typically having minimal details in a request item, so
	// the merging of multiple requests with the same url/method type.. could result in a more complete resource.
	Version string `json:"version"`

	// This is the OWNER of the resource origination. OpenAPI includes a Info section with Title in it. This would be used in combination with
	// Version to ensure a unique resource exists AND yet allows two (or more) API definition sources with the SAME title but different versions
	// to potentially be loaded while ensuring the latest version of the Resource is stored in the pool of resources.
	Owner string `json:"owner"`

	// This field can be set to TRUE (by default it should be) by a loader IF the currently processing operation/path/resource IS the latest version
	// in the case where two (or more) different versions of the same API are being loaded in one execution. Because the tool can load multiple
	// resources in one execution.. it is possible that multiple sources are of the same API, but different versions. This flag will allow
	// generators to work with ONLY the latest resources for something like code generation (to ensure only a single function for example).. but also
	// allow other generators such as documentation generators the ability to use ALL the versions of resources for any purpose such as generating
	// docs showing the different versions, or maybe a generator that compares differences in a graph or something.
	Latest bool `json:"latest"`
}

type QueryIn string

const (
	PATH   QueryIn = "path"
	HEADER QueryIn = "header"
	COOKIE QueryIn = "cookie"
	QUERY  QueryIn = "query"
)

type Parameter struct {
	Name              string     // The original json parameter name, eg param_name
	In                QueryIn    // Where the parameter is defined - path, header, cookie, query
	Description       string     // description of this parameter
	Required          bool       // Is this a required parameter
	Type              string     // the type of parameter (string, int, number, etc)
	Format            string     // The format of the Type property, e.g. int32, float64
	VariableNameValue string     // The name of variable for the value
	VariableNameKey   string     // The name of variable for the key
	Value             string     // A parameter can have a value
	Components        Components // Reference to the component created for this parameter if the type is object or array
}

// This describes a request body
type Request struct {
	Required    bool
	ContentType string
	// If this is a reference to a components/schema or components/requestBodies, store it
	Ref string
	// simpler type of content it is.  JSON, XML, etc. Easier to use programatically than parsing
	// application/json in some custom function or in templates.
	Type    string
	Default bool
	Schema  *Component
}

// This describes a request response
type Response struct {
	Status         string
	Description    string
	ResponseBodies ResponseBodies
}

type ResponseBody struct {
	MediaType string
	Ref       string
	Default   bool
	Schema    *Component
	Example   string
}

// MakeResourceName
// This function will return a created resource name as a string from the provided parameters. It can be used
// for creating names for resources even when the operationId is emtpy.
func MakeResourceName(name, method, path, subPath string) string {
	var resourceName string

	if name == "" || len(name) <= 0 {
		if subPath == "/" || len(subPath) <= 0 {
			resourceName = ToCamelCase(method+RemoveWhiteSpaceAndCaps(path), false)
		} else {
			resourceName = ToCamelCase(method+RemoveWhiteSpaceAndCaps(subPath), false)
		}
	} else {
		resourceName = ToCamelCase(RemoveWhiteSpaceAndCaps(name), false)
	}

	return resourceName
}

// This function can be used by resources to create the unique ResourceID field composed of parts like method and path.
//
// method should consist of a value of 'get' 'post' 'put' 'delete' 'patch' or 'option'
func (Resources) MakeResourceId(method, path string) string {
	// remove any {} and {{}}

	return method + ":" + path
}

func (r Resources) GetResourcesByHierarchy() map[string]*Resources {
	m := make(map[string]*Resources, 0)

	for _, resource := range r {
		if nil != resource {
			resources := m[resource.Root]

			if nil == resources {
				r := make(Resources, 0)
				m[resource.Root] = &r
				resources = m[resource.Root]
			}

			*resources = append(*resources, resource)
		}
	}

	return m
}

func (r Resources) GetLatestResources() *Resources {
	resources := make(Resources, 0)

	// if for some reason r is nil.. return an empty Resources object
	if nil == r {
		return &resources
	}

	for _, resource := range r {
		if resource.Latest {
			resources = append(resources, resource)
		}
	}

	return &resources
}

// MakeUniqueId
//
// This function will create a unique string value from the provided path.
// Basically it will convert any path parameter variables (denoted with {value}) to
// values without the { } wrapper, convert / to :
//
//	:/?#[]@!$&'()*+,;=  are url encoded values in a query string (after ?)
//
// ex: /users/{userId}/stores/{storeId}/purchases would become users:userId:stores:storeId:purchases
//
// query parameters are removed as they are dynamic to the path.. e.g. while path with variables can have
// dynamic values as part of the path, the query parameters are not part of the path resource identifier. They
// are provided to the resource as key/value pairs and don't identify the resource.
// resources: https://github.com/OAI/OpenAPI-Specification/issues/182
//
//	https://spec.openapis.org/oas/v3.1.0#patterned-fields
//	3.2+ : https://github.com/OAI/OpenAPI-Specification/issues/2572
func (r *Resources) MakeUniqueId(path, method string) string {
	s := path

	// remove anything past the ? if it exists..
	indx := strings.Index(path, "?")

	if indx > 0 {
		s = path[:indx]
	}

	// if path starts with /, remove it
	if len(path) > 0 && path[0] == '/' {
		s = s[1:]
	}

	pathVars := regexp.MustCompile(`{[^{}]*}`)
	matches := pathVars.FindAllStringSubmatch(s, -1)

	for _, ss := range matches {
		val := ss[0]
		v := val[1 : len(val)-1]
		s = strings.ReplaceAll(s, val, v)
	}

	return method + ":" + s
}

// NewResource
func (r *Resources) NewResource(path, method, name, description, summary, source, version, owner string, deprecated, latest bool, resourceType ResourceType) (resource *Resource, err error) {
	// make sure some must have fields are not nil/emtpy
	if len(path) <= 0 || len(method) <= 0 {
		err = fmt.Errorf(" resource must have a path and method ")
	} else {
		// Get the root path of the provided path value
		root := path

		if path[0] == '/' {
			pth2 := path[1:]
			indx2 := strings.Index(pth2, "/")

			if indx2 > 0 {
				root = path[:indx2+1]
			}
		}

		resource = &Resource{
			Id:           generateUniqueInt(),
			Path:         path,
			Method:       method,
			ResourceId:   r.MakeUniqueId(path, method),
			Root:         root,
			Name:         name,
			Description:  description,
			Summary:      summary,
			Source:       source,
			Version:      version,
			Owner:        owner,
			Deprecated:   deprecated,
			Latest:       latest,
			ResourceType: resourceType,
			Variables:    make(map[string]string, 0),
			Parameters:   make(Parameters, 0),
			Requests:     make(Requests, 0),
			Responses:    make(Responses, 0),
		}

		// auto set and increment the ID field
		*r = append(*r, resource)
	}

	return
}

func (r Resources) FindResourceByName(name string) *Resource {
	if nil == r {
		return nil
	}

	for _, res := range r {
		equal := strings.EqualFold(name, res.Name)
		if equal {
			return res
		}
	}

	return nil
}

// This function will utilise the reciever Resources object to try to find
// the resource based on the passed in resourceId. If found, returned.. nil otherwise
func (r Resources) FindResource(path, method string, owner *string) *Resource {
	if nil == r {
		return nil
	}

	id := r.MakeUniqueId(path, method)

	for _, res := range r {
		if res.ResourceId == id && res.Owner == *owner {
			return res
		}
	}

	return nil
}

func (r Resources) FindResourceByUuid(id int) *Resource {
	if nil == r {
		return nil
	}

	for _, res := range r {
		if res.Id == id {
			return res
		}
	}

	return nil
}

func (r *Resources) RemoveResource(resourceId string) {
	if nil == r {
		return
	}

	for pos, res := range *r {
		if res.ResourceId == resourceId {
			*r = append((*r)[:pos], (*r)[pos+1:]...)
		}
	}
}
