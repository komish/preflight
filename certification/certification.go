package certification

import "github.com/sirupsen/logrus"

// Check as an interface containing all methods necessary
// to use and identify a given check.
type Check interface {
	// Validate checks whether the asset enforces the check.
	Validate(image string, logger *logrus.Logger) (result bool, err error)
	// Name returns the name of the check.
	Name() string
	// Metadata returns the check's metadata.
	Metadata() Metadata
	// Help return the check's help text.
	Help() HelpText
	// IsBundleCheck returns whether this check is for a bundle or not
	IsBundleCheck() bool
	// IsBundleCompatible returns whether this check is suitable for a bundle
	IsBundleCompatible() bool
}

// Metadata contains useful information regarding the check
type Metadata struct {
	Description      string `json:"description" xml:"description"`
	Level            string `json:"level" xml:"level"`
	KnowledgeBaseURL string `json:"knowledge_base_url,omitempty" xml:"knowledgeBaseURL"`
	CheckURL         string `json:"check_url,omitempty" xml:"checkURL"`
}

// HelpText is the help message associated with any given check
type HelpText struct {
	Message    string `json:"message" xml:"message"`
	Suggestion string `json:"suggestion" xml:"suggestion"`
}

type CheckInfo struct {
	Metadata `json:"metadata" xml:"metadata"`
	HelpText `json:"helptext"`
}

func NewCheck() Check {
	return &DefaultCheck{}
}

type DefaultCheck struct{}

func (dc *DefaultCheck) Validate(image string, logger *logrus.Logger) (bool, error) {
	return false, nil
}

func (dc *DefaultCheck) Name() string {
	return ""
}

func (dc *DefaultCheck) Metadata() Metadata {
	return Metadata{}
}

func (dc *DefaultCheck) Help() HelpText {
	return HelpText{}
}

func (dc *DefaultCheck) IsBundleCheck() bool {
	return false
}

func (dc *DefaultCheck) IsBundleCompatible() bool {
	return true
}
