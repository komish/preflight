package certification

// Policy as an interface containing all methods necessary
// to use and identify a given policy.
type Policy interface {
	// Validate whether the asset enforces the policy.
	Validate(image string) (result bool, log []byte, err error)
	// return the name of the policy
	Name() string
	// return the policy's metadata
	Metadata() Metadata
	// return the policy's help text
	Help() HelpText
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

type Log = []byte

type PolicyWithRuntimeLog struct {
	Policy Policy
	Log    Log `json:"log,omitempty" xml:"log,omitempty"`
}

type PolicyInfo struct {
	Metadata `json:"metadata" xml:"metadata"`
	HelpText `json:"helptext"`
}

type PolicyMetadataWithLog struct {
	Metadata
	Log `json:"log,omitempty" xml:"log,omitempty"`
}

type PolicyInfoWithLog struct {
	PolicyInfo
	Log `json:"log,omitempty" xml:"log,omitempty"`
}

type PolicyHelpTextWithLog struct {
	HelpText
	Log `json:"log,omitempty" xml:"log,omitempty"`
}
