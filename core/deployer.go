package core

/*
type deployer struct {
	Deployer Deployer
	Name     string
}

type Deployers []*deployer

var deployers Deployers

func init() {
	// initialize the singleton instance of the map that contains all loaders
	deployers = make(Deployers, 0)
}

// This interface should be implicitly implemented by concrete SDK generators.
type Deployer interface {
	// Deploy
	// string parameter is the path to find artifacts to deploy
	// second param is for passing in a dynamic chunk of JSON or YAML that each implementation
	// can do whatever it needs to do with it. This could be to pass in authentication details,
	// deployment details, etc. Each deployer implementation can define it's own Struct type
	// and document that.. and the engine can fine/pass this info to the deployer if provided
	// at runtime in some manner.
	Deploy(string, map[string]interface{}) error
}

func GetDeployers() Deployers {
	return deployers
}

func (g *Deployers) New(name string, deplyr Deployer) {
	dep := &deployer{
		Name:     name,
		Deployer: deplyr,
	}

	deployers = append(deployers, dep)
	*g = deployers
}

// FindGeneratorByName
//
// This receiver function will attempt to find a generator that has a name-variant string value that matches the name string parameter
// provided by looping through all registered generators. If a match is found, it is returned otherwise nil is returned indicating no match.
func (d Deployers) FindDeployerByName(name string) Deployer {
	for _, dep := range d {
		if name == dep.Name {
			return dep.Deployer
		}
	}

	return nil
}


*/
