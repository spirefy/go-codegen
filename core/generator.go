package core

import "github.com/spirefy/codegen/types"

// This object wraps the details of an individual generator implementation
//
// Generator is the pointer to an implementation instance. Ideally only a single instance is created for use
// Name is the name of the generator and represents the way to match the generator(s) to use at runtime
// TemplateFS if provided is a pointer to embeded templates.
// TemplatePath is a path provided as the annotation to the go:embed associated with the TemplateFS provided
// RootPath is the root directory.. temporary or otherwise, where generated output will be written to
// Deployer if provided is the name of deployer implementation to use to deploy with.
type generator struct {
	Name      string
	Variant   string
	Type      string
	Options   any
	Generator Generator
}

type Generators []*generator

var generators Generators

// Generator
//
// This interface should be implicitly implemented by concrete SDK generators.
type Generator interface {
	// Generate
	// This function is called to generate output.
	// The first parameter is the output path to generate to.
	// The second parameter is the pointer to the Resources loaded prior to this call.. representing the items to generate output from
	// The third parameter is a slice of Unique components. These will typically be request or response payload components, but could be derived
	//   however each loader needs. Loaders should look to see if a component exists before adding one.
	// The fourth parameter is a slice of workflows.. each representing a potential series of procedures and steps that reference resources (by uuid) and include
	//   details on response values to pull out in to variables and request variables as well.
	Generate(string, GeneratorContext) error
}

type TypeGenerator interface {
	// Generate
	// This method will generate the specified type to the specified output path for the provided Components
	// First parameter is the Components structure. This is the combined set of all loaded components across all
	// sources... same as provided to generators
	// Second parameter is the target matching that provided by the comma separated list of targets
	// Third parameter is the output path where the targets generated code will go
	Generate(types.Components, string, string)
}

func init() {
	// initialize the singleton instance of the map that contains all loaders
	generators = make(Generators, 0)
}

// GetGenerators
//
// This function returns the singleton instance of the generator map... type Generators. This map
// contains 0 or more generator implementations. Each generator implementation is a singleton instance.
// This is the function to call by generator implementations to register a generator implementation...
// example:
//
// // init function called before other funcs in codegen core.. to ensure any imported generator implementation
// // is loaded first.. so it is ready for use.
//
//	func init() {
//	  generators := GetGenerators()
//	  generators.New(name, variant, typ, GeneratorImpl)
//	}
func GetGenerators() Generators {
	return generators
}

func (g *Generators) New(name, variant, typ string, genertr Generator) {
	gen := &generator{
		Name:      name,
		Variant:   variant,
		Type:      typ,
		Generator: genertr,
	}

	generators = append(generators, gen)
	*g = generators
}

// FindGeneratorByName
//
// This receiver function will attempt to find a generator that has a name-variant string value that matches the name string parameter
// provided by looping through all registered generators. If a match is found, it is returned otherwise nil is returned indicating no match.
func (g *Generators) FindGeneratorByName(name, variant string) *generator {
	searchName := name + "-" + variant

	for _, gen := range *g {
		if searchName == gen.Name+"-"+gen.Variant {
			return gen
		}
	}

	return nil
}

/*
// LoadOptions
//
// This receiver function will attempt to load embedded options of the generator implementation and if found
// create the array of Options and set it on the generator Options property.
func (g generator) LoadOptions(embeded Embeded) error {
	if nil != embeded.EmbededFS && len(embeded.EmbededPath) > 0 {
		optionFile, _ := embeded.EmbededFS.ReadFile(embeded.EmbededPath[0])

		// Options, if provided are in the json format of
		// [
		//   {
		//     "name": "<name value>",
		//     "description": "<description value>",
		//     "type": "<type value (string, bool, int, float, etc)>",
		//     "value": "<actual value.. usually a default value to use initially but overwriteable at runtime>"
		//   }
		// ]
		//
		// So.. we use []interface{} instead of more commonly seen map[string]interface{} as we know the structure is always
		// an array of objects
		options := make([]interface{}, 0)
		err := json.Unmarshal(optionFile, &options)
		if err != nil {
			CLog.Log(CLogCore, "", "There was an error trying to unmarshal the options for generator: ", err)
			return err
		}

		// If we have options.. iterate the slice, converting each chunk in to a valid Option object, building the slice of Option objects, and assigning it to the generateObj.
		if nil != options {
			opts := make(Options, 0)

			for _, opt := range options {
				if nil != opt {
					jsonString, _ := json.Marshal(opt)
					o := Option{}
					err = json.Unmarshal(jsonString, &o)

					if err != nil {
						CLog.Log(CLogCore, "", "There was an error trying to unmarshal the options for generator: ", err)
						return err
					}

					opts = append(opts, &o)
				}
			}

			g.Options = opts
		}
	}

	return nil
}

*/
