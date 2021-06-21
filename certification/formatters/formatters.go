package formatters

import (
	"github.com/komish/preflight/certification/runtime"
)

// FormatterFunc describes a function that formats the policy validation
// results.
type FormatterFunc = func(runtime.Results) (response []byte, formatError error)

// func GetResponse(r runtime.Results) runtime.UserResponse {
// passedPolicies := make([]policy.Metadata, len(r.Passed))
// failedPolicies := make([]policy.Policy, len(r.Failed))
// erroredPolicies := make([]policy.Policy, len(r.Errors))

// if len(r.Passed) > 0 {
// 	for i, policyData := range r.Passed {
// 		passedPolicies[i] = policyData.GetMetadata()
// 	}
// }

// if len(r.Failed) > 0 {
// 	for i, policyData := range r.Failed {
// 		failedPolicies[i] = policy.Policy{
// 			Metadata: policyData.GetMetadata(),
// 			HelpText: policyData.GetHelp(),
// 		}
// 	}
// }

// if len(r.Errors) > 0 {
// 	for i, policyData := range r.Errors {
// 		erroredPolicies[i] = policyData.GetHelp()
// 	}
// }

// response := runtime.UserResponse{
// 	Image:             r.TestedImage,
// 	ValidationVersion: version.Version,
// 	Results: runtime.UserResponseText{
// 		Passed: passedPolicies,
// 		Failed: failedPolicies,
// 		Errors: erroredPolicies,
// 	},
// }

// return response

// }

// GenericJSONFormatter is a FormatterFunc that formats results as JSON
// func GenericJSONFormatter(r runtime.Results) ([]byte, error) {
// 	response := GetResponse(r)

// 	responseJSON, err := json.MarshalIndent(response, "", "    ")
// 	if err != nil {
// 		e := fmt.Errorf("%w with formatter %s: %s",
// 			errors.ErrFormattingResults,
// 			"json",
// 			err,
// 		)

// 		return nil, e
// 	}

// 	return responseJSON, nil
// }

// // GenericXMLFormatter is a FormatterFunc that formats results as XML
// func GenericXMLFormatter(r runtime.Results) ([]byte, error) {
// 	response := GetResponse(r)

// 	responseJSON, err := xml.MarshalIndent(response, "", "    ")
// 	if err != nil {
// 		e := fmt.Errorf("%w with formatter %s: %s",
// 			errors.ErrFormattingResults,
// 			"json",
// 			err,
// 		)

// 		return nil, e
// 	}

// 	return responseJSON, nil
// }

// func JUnitXMLFormatter(r runtime.Results) ([]byte, error) {
// 	return nil, fmt.Errorf("%w: The JUnit XML Formatter is not implemented", errors.ErrFeatureNotImplemented)
// }
