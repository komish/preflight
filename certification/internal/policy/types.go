package policy

// Policy as an interface containing all methods necessary
// to use and identify a given policy.
type Policy interface {
	// Validate whether the asset enforces the policy.
	Validate(image string) (result bool, err error)
	// return the name of the policy
	GetName() string
	// return the policy's metadata
	GetMetadata() Metadata
	// return the policy's help text
	GetHelp() HelpText
}

// Metadata contains useful information regarding the policy
type Metadata struct {
	Description      string `json:"description" xml:"description"`
	Level            string `json:"level" xml:"level"`
	KnowledgeBaseURL string `json:"knowledge_base_url,omitempty" xml:"knowledgeBaseURL"`
	PolicyURL        string `json:"policy_url,omitempty" xml:"policyURL"`
}

// HelpText is the help message associated with any given policy
type HelpText struct {
	Message    string `json:"message" xml:"message"`
	Suggestion string `json:"suggestion" xml:"suggestion"`
}
