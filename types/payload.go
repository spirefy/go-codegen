package types

type LoadedResponse struct {
	Resources  Resources  `json:"resources"`
	Components Components `json:"components"`
	Workflows  Workflows  `json:"workflows"`
}
