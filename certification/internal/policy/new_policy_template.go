package policy

import "github.com/komish/preflight/certification/errors"

// +preflight:codegen:replace-template-with-actual-name
func TemplatePolicy() *Definition {
	return &Definition{
		ValidatorFunc: isTemplateValidatorFunc,
		Metadata:      isTemplatePolicyMeta,
		HelpText:      isTemplatePolicyHelp,
	}
}

// +preflight:codegen:replace-template-with-actual-name
// +preflight:codegen:todos-prompt
var isTemplateValidatorFunc = func(image string) (bool, error) {
	// TODO implement validation logic here and change return value
	return false, errors.ErrFeatureNotImplemented
}

// +preflight:codegen:replace-template-with-actual-name
// +preflight:codegen:todos-prompt
var isTemplatePolicyMeta = Metadata{
	Description:      "TODO description here",
	Level:            "TODO specify level here",
	KnowledgeBaseURL: "TODO provide KB url",
	PolicyURL:        "TODO provide policy url",
}

// +preflight:codegen:replace-template-with-actual-name
// +preflight:codegen:todos-prompt
var isTemplatePolicyHelp = HelpText{
	Message:    "TODO policy message here",
	Suggestion: "TODO suggestion for meeting the policy criteria here",
}
