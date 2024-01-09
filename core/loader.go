package core

type loader struct {
	Name   string
	Type   string
	Loader Loader
}

type Loaders []*loader

// private instance of Loaders type to contain all loaders
var loaders Loaders

type Loader interface {
	// Load
	//
	// Allows the loading/parsing of a source entity in to a LoaderResources object. Upon return the codegen
	// tool will merge/combine/de-dup the LoaderResources response with the unified in memory model
	Load(string, Shared) (LoaderResources, error)
}

func init() {
	// initialize the loaders map
	loaders = make([]*loader, 0)
}

// GetLoaders
//
// This function can be used by consumers to retrieve the map of loaders.
func GetLoaders() Loaders {
	return loaders
}

// NewLoader
//
// This receiver function will add a new loader struct to the map of of loaders attached to the Loaders receiver parameter.
// It WILL replace a loader by the same name if it exists.
func (l *Loaders) NewLoader(name, typ string, ld Loader) {
	ldr := &loader{
		Name:   name,
		Type:   typ,
		Loader: ld,
	}

	loaders = append(loaders, ldr)
	*l = loaders
}

// FindLoaderByName
//
// This receiver function will return either a loader struct instance if the provided name parameter matches a key in the map, or nil
// otherwise.
func (l Loaders) FindLoaderByName(name string) *loader {
	for _, ldr := range l {
		if name == ldr.Name {
			return ldr
		}
	}

	return nil
}

// FindLoaderByType
//
// This receiver function will return a loader struct instance if a type is matched to the provided type value, otherwise nil
func (l Loaders) FindLoaderByType(typ string) *loader {
	for _, ldr := range l {
		if ldr.Type == typ {
			return ldr
		}
	}

	return nil
}
