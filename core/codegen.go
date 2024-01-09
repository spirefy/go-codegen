package core

import (
	"encoding/json"
	"fmt"
	"github.com/spirefy/codegen/types"
	pluginengine "github.com/spirefy/go-plugin-engine"
	pet "github.com/spirefy/go-plugin-engine/types"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

type SourceType int8

const ( // iota is reset to 0
	SourceFileType SourceType = iota
	SourceUrlType  SourceType = iota
	SourceAPIType  SourceType = iota
)

type (
	Aliases map[string]string

	// LoaderResources
	//
	// The primary resources (endpoints, components and workflows) contributed to by loaders during the loading/parsing process of sources
	LoaderResources struct {
		// Generic shared Resources object that all loaders will load in to and generators would use to generate code from
		Resources *types.Resources `json:"resources"`

		// Generic shared Components (payloads) referenced by Resources and would be loaded by loaders and used by generators to generate
		// static types for example.
		Components *types.Components `json:"components"`

		// Workflows represent a series of chainable Resources that provide details as to what is required by a given step (http request typically)
		// and what values to pull out of that steps response (http response typically). Each workflow has a series of steps, which each step would
		// indicate the Resource it depends upon (e.g. represents) to find referenced payload and parameter details. Code generators can then use
		// this workflow to assembled more bespoke code, such as robust SDKs that provide multiple API calls in a single function.
		Workflows *types.Workflows `json:"workflows"`
	}

	Shared map[string]any

	GeneratorContext struct {
		types.Resources  `json:"resources" yaml:"resources"`
		types.Components `json:"components" yaml:"components"`
		types.Workflows  `json:"workflows" yaml:"workflows"`

		// Any sort of variables that might be usable within the generator engine, loaders, deployers or generators
		Variables map[string]any `json:"variables" yaml:"variables"`
	}

	// CodegenConfig
	// Configuration structure for codegen core
	// Variables are for generators
	// Aliases are for loaders
	// BespokeWorkflow indicates reduction of resources and components found within the workflows. In other words, it
	//     will loop through all workflows and figure out which resources AND components are being referenced/used BY
	//     workflows and reduce the overall resources and workflows objects to only those that are being used by
	//     workflows. This should result in bespoke generated output matching only what workflows use and not the full
	//     set of resources and components referenced by loaded sources.
	CodegenConfig struct {
		// Any sort of variables that might be usable within the generator engine, loaders, deployers or generators
		Variables map[string]any `json:"variables" yaml:"variables"`

		// Aliases to replace if needed such as long URLs to short values
		Aliases map[string]string `json:"aliases" yaml:"aliases"`

		// Utilize bespoke output vs verbose output. Bespoke output should scan the workflow steps, determine what API
		// endpoints and components will be used including any reference trees, and build up a subset of the overall
		// component and resources loaded. This subset of resources and components is what the generator would then use
		// to generate output from. This is specifically related to workflows. Generators that are only generating
		// from resources and components would not make use of this flag.
		BespokeWorkflow bool `json:"bespoke" yaml:"bespoke"`

		// Lint will attempt to lint/validate sources depending on options. If set to true, apply lint.
		// TODO: Add support for individual source type linters when supported
		Lint bool `json:"lint" yaml:"lint"`

		// Validate will determine how linting is used. Validate set to false, will allow logging of linting details
		// but continue to generate artifacts even if the result of linting would be an invalid source. True will
		// indicate that if linting results in any sort of invalid source, that particular source should not be used
		// when loading the series of sources in to the unified structure.
		Validate bool `json:"validate" yaml:"validate"`
	}

	Sources struct {
		Type   SourceType `json:"type" yaml:"type"`
		Source []byte     `json:"source" yaml:"source"`
		Path   string     `json:"path" yaml:"path"`
		Loaded bool
	}

	Targets struct {
		Generator Generator `json:"generator" yaml:"generator"`
		Name      string    `json:"name" yaml:"name"`
		Variant   string    `json:"variant" yaml:"variant"`
		Type      string    `json:"type" yaml:"type"`
		// Options   Options   `json:"options" yaml:"options"`
		// Configuration that contains variables such as those passed on the CLI or as part of any API service request..
		// as well as aliases that can be used to reduce long values in to shorter ones for code generation for example.
		// Per Generator configuration
		Configuration *CodegenConfig `json:"configuration" yaml:"configuration"`
	}
)

type Codegen struct {
	sync.Mutex   // Add a mutex for synchronization
	PluginEngine *pluginengine.Engine
	Components   *types.Components
	Resources    *types.Resources
	Workflows    *types.Workflows
}

// codegen is a single instance of the Codegen engine, returned by the NewCodegen func and instantiated if not already.
var (
	codegen *Codegen
	once    sync.Once
)

func (a Aliases) FindAlias(key string) *string {
	for k, v := range a {
		if k == key {
			return &v
		}
	}

	return nil
}

// LoadSourceContents
//
// This function is a helper function to load contents found at the source location.
// If the source starts with http:// or https:// an attempt to download the contents using
// a http GET request is made. Otherwise an attempt to open a local file using the source
// string as the path and filename/extension is attempted to load.
func (c *Codegen) LoadSourceContents(source string) ([]byte, error) {
	var reader io.Reader
	var err error

	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		var resp *http.Response
		resp, err = http.DefaultClient.Get(source)

		if err != nil {
			return nil, err
		}

		reader = resp.Body
	} else {
		var file *os.File
		file, err = os.Open(source)

		defer func(r *os.File) {
			err := r.Close()
			if err != nil {
				CLog.Error(CLogCore, "", "Unable to load source: %s"+source)
			}
		}(file)

		if err != nil {
			return nil, err
		}

		// assign the io.Reader to the os.File
		reader = file
	}

	contents, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (c *Codegen) Load(source Sources, shared Shared) error {
	// Lock so that this call can operate in a thread safe manner.
	c.Lock()
	defer c.Unlock()

	// loop through ALL registered and try to load the source...
	// if the response from a loader is nil.. it loaded and ignore asking any other
	// loaders.. but if there is an error.. log error and try the next one in the list
	for _, ldr := range loaders {
		switch source.Type {
		case SourceUrlType:
			lrs, err := ldr.Loader.Load(source.Path, shared)

			if err != nil {
				CLog.Error(CLogLoaders, source.Path, "Unable to load source: %s", source)
				continue
			}

			fmt.Println("Doing something with LRS from URL: ", lrs)
		case SourceFileType:
			lrs, err := ldr.Loader.Load(source.Path, shared)
			// TODO: Check for "event" extension point from plugins for any plugin listening for LOADED events

			if err != nil {
				CLog.Error(CLogLoaders, source.Path, "Unable to load source: %s", source)
				continue
			}
			fmt.Println("Doing something with LRS from PATH: ", lrs)
		case SourceAPIType:
			// err := ldr.Loader.Load(src., loaderContext, lntr)
			//
			// if err != nil {
			// 	return fmt.Errorf("error loading local file %s", err)
			// }
		}
	}

	return nil
}

// Generate
//
// # This function will be called by a consumer of the library
func (c *Codegen) Generate(sources []Sources, generators []Targets, outputPath string, typeGenerator TypeGenerator, config *CodegenConfig) error {
	shared := Shared{}
	shared["aliases"] = config.Aliases

	loaders := GetLoaders()

	// If we have loaders AND source files...
	if len(loaders) > 0 && len(sources) > 0 {
		// Split the sources by , to get each source

		// Now loop through all loader types, and split the source files it contains
		for _, src := range sources {
			err := c.Load(src, shared)

			if err != nil {
				CLog.Log(CLogCore, "", "source loading error: ", err)
			}
		}
	}

	// Now that loading is done lets see if we need to reduce the resources and components if the bespoke option
	// is provided and workflows exist
	// Workflows and steps, and identifying ONLY the resources/components used within the workflows provided. This
	// will reduce the overall SDK code generated to only what is used.
	if config.BespokeWorkflow && nil != c.Workflows && len(*c.Workflows) > 0 {
		// comps := make(Components, 0)
		rescs := make(types.Resources, 0)

		// TODO: FIX THIS
		/*

			// loop through workflows
			for _, w := range *loaderContext.Workflows {
				// loop through all steps of this workflow
				for _, s := range w.Steps {
					// check if the step resource is already in the rescs slice
					alreadyAdded := false
					for _, r := range rescs {
						if r == s.Resource {
							alreadyAdded = true
							break
						}
					}

					// not yet added, lets do so now
					if !alreadyAdded {
						// TODO: REMOVE THIS ONCE ITS ALL WORKING PROPERLY (AND A TEST IS ADDED FOR THIS)
						if nil != s.Resource {
							rescs = append(rescs, s.Resource)
						}

						// OK.. adding this resource.. so now look for any/all components this resource might use/reference
						// and build up the list of bespoke components
						// TODO: Need to add this code for bespoke components
					}
				}
			}

		*/

		fmt.Println("Resources len is now: ", len(rescs))
		// replace ALL resources with just the bespoke resources
	}

	// As long as we have resources & generators...
	if nil != *c.Resources && len(*c.Resources) > 0 && nil != generators && len(generators) > 0 {
		for _, gen := range generators {
			if nil != gen.Generator {
				// if configuration variables exist, lets apply them to the generator object automatically setting
				// any json/yaml variables that match json/yaml tags of any type defined for the generator
				if nil != config && nil != config.Variables {
					jsonString, err := json.Marshal(config.Variables)

					if nil != err {
						CLog.Log(CLogCore, "", "Error converting variables to json: ", err)
					} else {
						err = json.Unmarshal(jsonString, gen.Generator)

						if nil != err {
							CLog.Log(CLogCore, "", "Error trying to unmarshal variable string to generator object: %v \n\n variables: %v ", err, jsonString)
						}
					}
				}

				// Adjust the output path to include the generator name (lower cased)
				// TODO: Make this a default way to determine output, but optional output details may be provided
				// TODO: via config perhaps.
				s := string(os.PathSeparator)
				pth := strings.TrimSuffix(outputPath, s)
				var lastPart string

				if len(gen.Type) > 0 {
					lastPart = gen.Name + s + gen.Type + s
				} else {
					lastPart = gen.Name + s
				}

				if len(gen.Variant) > 0 {
					lastPart = lastPart + gen.Variant + s
				}

				pth = pth + s + strings.ToLower(lastPart)

				// if the type generator is available, pass the components, target generator and output path
				if nil != typeGenerator {
					typeGenerator.Generate(*c.Components, gen.Name+"-"+gen.Variant, pth)
				}

				generatorContext := GeneratorContext{
					Resources:  *c.Resources,
					Components: *c.Components,
					Workflows:  *c.Workflows,
					Variables:  config.Variables,
				}

				err := gen.Generator.Generate(pth, generatorContext)
				if nil != err {
					CLog.Log(CLogCore, "", "Error generating: ", err)
				}
			}
		}
	}

	return nil
}

type LoadedResponse struct {
	Resources  types.Resources  `json:"resources"`
	Components types.Components `json:"components"`
	Workflows  types.Workflows  `json:"workflows"`
}

// LoaderExtensions
// This function will be called with all extensions related to the codegen.plugins.loaders extension point. It will
// loop through each extension, calling an exepcted external plugin func called GetLoaderDetails. This function will
// return a json payload conforming to the following json structure:
//
//		{
//	   "resources": [...],
//	   "components": [...],
//	   "workflows": [...]
//		}
func LoaderExtensions(extensions []*pluginengine.Extension) error {
	for _, ex := range extensions {
		_, data, err := ex.Plugin.Call(ex.FuncName, nil)
		if nil != err {
			fmt.Println("Error calling extension func: ", err)
		}
		lr := LoadedResponse{}
		err = json.Unmarshal(data, &lr)
		if err != nil {
			fmt.Println("ERR unmarshalling: ", err)
		}

		if nil != lr.Resources && len(lr.Resources) > 0 {
			for _, r := range lr.Resources {
				fmt.Println("Resource: ", r)
			}
		}
	}

	return nil
}

func GeneratorExtensions(extensions []*pluginengine.Extension) error {
	for _, ex := range extensions {
		_, data, err := ex.Plugin.Call(ex.FuncName, nil)
		if nil != err {
			fmt.Println("Error calling extension func")
		}

		fmt.Println("Data: ", data)
	}

	return nil
}

// NewCodegen
//
// This function is called by an application incorporating the codegen engine into it. It will create a new instance
// of the codegen engine. Typically there wouldn't be a reason to create more than one instance so this uses a local
// variable (not exported) to ensure a static singleton is returned

func NewCodegen(path string) *Codegen {
	once.Do(func() {
		resources := make(types.Resources, 0)
		workflows := make(types.Workflows, 0)
		components := make(types.Components, 0)
		engine := pluginengine.NewPluginEngine()

		codegen = &Codegen{
			Mutex:        sync.Mutex{},
			PluginEngine: engine,
			Components:   &components,
			Resources:    &resources,
			Workflows:    &workflows,
		}
	})

	// Set up plugin engine HOST functions for this codegen engine for loaders, generators. Host application can add
	// additional extension points as desired for its specific use of the codegen engine.
	loaderExtensionPoint := pluginengine.ExtensionPoint{
		ExtensionPoint: pet.ExtensionPoint{
			Id:          "codegen.plugins.loaders",
			Description: "This extension point allows plugins to provide extensions which implement loading/parsing a source format and return the parsed results as a json structure that matches the schema",
			Name:        "Codegen Source Loader Plugin",
			StartOnLoad: true,
		},
		Func: LoaderExtensions,
	}

	generatorExtensionPoint := pluginengine.ExtensionPoint{
		ExtensionPoint: pet.ExtensionPoint{
			Id:          "codegen.plugins.generators",
			Description: "This extension point allows plugins to provide extensions which implement generation of output from the provided codegen unified model of resources, components and workflows.",
			Name:        "Codegen Source Loader Plugin",
			StartOnLoad: true,
		},
		Func: GeneratorExtensions,
	}

	codegen.PluginEngine.RegisterHostExtensionPoint(loaderExtensionPoint)
	codegen.PluginEngine.RegisterHostExtensionPoint(generatorExtensionPoint)

	// load plugins
	err := codegen.PluginEngine.Load(path)
	if nil != err {
		fmt.Println("Error loading plugin engine: ", err)
	}

	return codegen
}
