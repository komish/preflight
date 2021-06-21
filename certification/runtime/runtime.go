package runtime

import (
	"github.com/komish/preflight/certification/internal/policy"
	"github.com/komish/preflight/version"
)

type Config struct {
	Image           string
	EnabledPolicies []string
}

type PolicyRunner struct {
	Image    string
	Policies []policy.Policy
	Results  Results
}
type Results struct {
	TestedImage string
	Passed      []policy.Policy
	Failed      []policy.Policy
	Errors      []policy.Policy
}

type UserResponse struct {
	Image             string                 `json:"image" xml:"image"`
	ValidationVersion version.VersionContext `json:"validation_lib_version" xml:"validationLibVersion"`
	Results           UserResponseText       `json:"results" xml:"results"`
}

type UserResponseText struct {
	Passed []policy.Policy
	Failed []policy.Policy
	Errors []policy.Policy
	// TODO: Errors does not actually include any error information
	// and it needs to do so.
}
