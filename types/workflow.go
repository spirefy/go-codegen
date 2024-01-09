package types

import "log"

// Input is a structure making up expected input parameters (variables) to Operations. They would be fed in from consumers of generated code, for example. They might appear as
// part of a function or method signature that are provided by a consumer of the generated code and the values would be derived outside of the generated code the workflow might
// generate.
type Input struct {
	Id string `json:"id"`
}

// This represents a single expression. The purpose of an expression is to define logic to be applied at runtime (via generated output)
// to derive at the final value that is the output that this expression is referenced for
type Expression struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

// This is the structure of an individual output
type Output struct {
	Id         string     `json:"id"`
	Expression Expression `json:"expression"`
}

type WorkflowParameter struct {
	Name   string `json:"name" yaml:"name"`
	In     string `json:"in" yaml:"in"`
	Value  string `json:"value" yaml:"value"`
	Style  string `json:"style,omitempty" yaml:"style,omitempty"`
	Target string `json:"target,omitempty" yaml:"target,omitempty"`
}

type Action struct {
	// Unique string representing this Action. This should be unique across the entire workflow.
	Id string `json:"id"`
}

// Step
// A single step object..
// Every step MUST reference a Resource OR another Step.. they are mutually exclusive. While both indicate omitempty, validation should be implemented to ensure one or the other is present
// Parameters specify 0 or more values to pass to this step.
// DependsOn if provided requires that any listed step has executed prior to this step.. to ensure an order is exeuctued and parallel execution of steps is not performed. Primarily for ensuring that parameters required by this step are available.
// Success
type Step struct {
	// Unique string representing this step. This should be unique across the entire workflow even though steps are part of operations.
	Id string `json:"id"`

	// A name if provided by the workflow source for this step.
	Name string `json:"name,omitempty"`

	// A description of what this step does. CommonMark may be used for rich text representation
	Description string `json:"description,omitempty"`

	// A reference to a Resource that this Step relates to. This
	Resource *Resource `json:"resource,omitempty"`

	// A map representing parameters to pass to an operation as specified in the references Resource parameters.
	Parameters WorkflowParameters `json:"parameters,omitempty"`

	// A slice of steps that MUST be completed sequentially before this Step can being execution.
	DependsOn []Step `json:"dependson,omitempty"`

	// Outputs is a map of Output objects, keyed on a friendly name and a value determined at runtime. The value can be an expression.
	Outputs map[string]Output `json:"outputs"`

	// SuccessCriteria is a slice of expressions that should be evaulated at rutnime (e.g. generated code would execute at runtime) to determine the success or failure of the step.
	SuccessCriteria []Expression `json:"successCriteria"`

	// Slice of actions that are executed if the SuccessCriteria is deemed successful (2xx StatusCode ??)
	OnSuccess []Action `json:"onSuccess,omitempty"`

	// Slice of actions that are executed if the SuccessCriteria is deemed failure (4xx/5xx StatusCode ???)
	OnFailure []Action `json:"onFailure,omitempty"`
}

type Steps []*Step

type Workflow struct {
	// A unique id for this operation... unique across the entire workflow
	Id          string             `json:"id"`
	Description string             `json:"description,omitempty"`
	Inputs      Components         `json:"inputs"`
	Steps       Steps              `json:"steps"`
	Outputs     map[string]*Output `json:"outputs"`
}

type (
	Workflows          []*Workflow
	WorkflowParameters map[string]WorkflowParameter
)

// AddWorkflow
//
// This method will add a workflow to the slice of Workflows provided by receiver *wf
// if the id provided does not match an existing id.. keeping the slice to unique workflows.
// It will return nil if the workflow exists.. and the workflow instance if it was added.
func (wf *Workflows) AddWorkflow(workflow *Workflow) {
	if nil != workflow && len(workflow.Id) > 0 {
		// make sure workflow ID does not already exist.. keeping workflows unique
		for _, w := range *wf {
			if w.Id == workflow.Id {
				log.Printf("Workflow with id %s already exists and can not be added", workflow.Id)
				return
			}
		}

		*wf = append(*wf, workflow)
	}
}
