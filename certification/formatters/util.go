package formatters

import (
	"github.com/komish/preflight/certification"
	"github.com/komish/preflight/certification/runtime"
	"github.com/komish/preflight/version"
)

// getResponse will extract the runtime's results and format it to fit the
// UserResponse definition in a way that can then be formatted.
func getResponse(r runtime.Results) runtime.UserResponse {
	passedPolicies := make([]certification.PolicyMetadataWithLog, len(r.Passed))
	failedPolicies := make([]certification.PolicyInfoWithLog, len(r.Failed))
	erroredPolicies := make([]certification.PolicyHelpTextWithLog, len(r.Errors))

	if len(r.Passed) > 0 {
		for i, execution := range r.Passed {
			passedPolicies[i] = certification.PolicyMetadataWithLog{
				Metadata: execution.Policy.Metadata(),
				// Log:      execution.Log,
			}
		}
	}

	if len(r.Failed) > 0 {
		for i, execution := range r.Failed {
			failedPolicies[i] = certification.PolicyInfoWithLog{
				PolicyInfo: certification.PolicyInfo{
					Metadata: execution.Policy.Metadata(),
					HelpText: execution.Policy.Help(),
				},
				// Log: execution.Log,
			}
		}
	}

	if len(r.Errors) > 0 {
		for i, execution := range r.Errors {
			erroredPolicies[i] = certification.PolicyHelpTextWithLog{
				HelpText: execution.Policy.Help(),
				// Log:      execution.Log,
			}
		}
	}

	response := runtime.UserResponse{
		Image:             r.TestedImage,
		ValidationVersion: version.Version,
		Results: runtime.UserResponseText{
			Passed: passedPolicies,
			Failed: failedPolicies,
			Errors: erroredPolicies,
		},
	}

	return response
}
