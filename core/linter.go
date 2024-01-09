package core

/*
import "fmt"

// linter
// This struct will provide a pluggable linter capability to allow sources to be linted. Each linter implementation ca
// work with specific source types, such as OpenAPI, WorkflowSpec, Collection, etc. This allows validation of input
// sources to ensure they pass any rulesets before any sort of artifact generation occurs.
type linter struct {
	Name    string
	Type    string
	Options *Options
	Linter  Linter
}

type Linters []*linter

var linters Linters

// Linter
//
// This interface should be implicitly implemented by concrete Linter implementations.
type Linter interface {
	// The first parameter is the
	Lint(string, []byte) error
}

func init() {
	// initialize the singleton instance of the map that contains all loaders
	linters = make(Linters, 0)
}

// GetLinters
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
func GetLinters() Linters {
	return linters
}

func (l *Linters) New(name, typ string, lnt Linter) {
	lint := &linter{
		Name:   name,
		Type:   typ,
		Linter: lnt,
	}

	linters = append(linters, lint)
	*l = linters
}

// FindLinterByName
//
// This receiver function will attempt to find a generator that has a name-variant string value that matches the name string parameter
// provided by looping through all registered generators. If a match is found, it is returned otherwise nil is returned indicating no match.
func (l Linters) FindLinterByName(name string) *linter {
	for _, lint := range l {
		if (name == lint.Name) || (name == lint.Name+"-"+lint.Type) {
			return lint
		}
	}

	return nil
}

func (l Linters) FindLinterByType(typ string) *linter {
	fmt.Println("Attempting to find linter by type: ", typ)

	for _, lint := range l {
		fmt.Println("Looking at linter named: ", lint.Name)
		if typ == lint.Type {
			return lint
		}
	}

	return nil
}

*/
