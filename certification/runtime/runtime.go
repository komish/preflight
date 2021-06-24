package runtime

import (
	"github.com/komish/preflight/certification"
	"github.com/komish/preflight/version"
)

type Config struct {
	Image           string
	EnabledPolicies []string
	ResponseFormat  string
}
type Results struct {
	TestedImage string
	Passed      []certification.PolicyWithRuntimeLog
	Failed      []certification.PolicyWithRuntimeLog
	Errors      []certification.PolicyWithRuntimeLog
}

type UserResponse struct {
	Image             string                 `json:"image" xml:"image"`
	ValidationVersion version.VersionContext `json:"validation_lib_version" xml:"validationLibVersion"`
	Results           UserResponseText       `json:"results" xml:"results"`
}

type UserResponseText struct {
	Passed []certification.PolicyMetadataWithLog
	Failed []certification.PolicyInfoWithLog
	Errors []certification.PolicyHelpTextWithLog
	// TODO: Errors does not actually include any error information
	// and it needs to do so.
}
