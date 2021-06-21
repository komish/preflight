package runtime

import (
	"fmt"

	"github.com/komish/preflight/certification/errors"
	"github.com/komish/preflight/certification/internal/policy"
)

// TODO: Decide what of this file actually needs exporting

type PolicyEngine interface {
	ExecutePolicies()
	// StorePolicies(...[]certification.Policy)
	GetResults() Results
}

// Register all policies
var NameToPoliciesMap = map[string]policy.Policy{
	policy.RunAsNonRootPolicy{}.GetName():     policy.RunAsNonRootPolicy{},
	policy.UnderLayerMaxPolicy{}.GetName():    policy.UnderLayerMaxPolicy{},
	policy.HasRequiredLabelPolicy{}.GetName(): policy.HasRequiredLabelPolicy{},
	policy.BasedOnUbiPolicy{}.GetName():       policy.BasedOnUbiPolicy{},
}

func GetPoliciesByName() []string {
	all := make([]string, len(NameToPoliciesMap))
	i := 0

	for k := range NameToPoliciesMap {
		all[i] = k
		i++
	}
	return all
}

func NewForConfig(config Config) (*PolicyRunner, error) {
	if len(config.EnabledPolicies) == 0 {
		// refuse to run if the user has not specified any policies
		return nil, errors.ErrNoPoliciesEnabled
	}

	policies := make([]policy.Policy, len(config.EnabledPolicies))
	for i, policyString := range config.EnabledPolicies {
		// search policies by names
		policy, exists := NameToPoliciesMap[policyString]
		if !exists {
			err := fmt.Errorf("%w: %s",
				errors.ErrRequestedPolicyNotFound,
				policyString)
			return nil, err
		}
		policies[i] = policy
	}
	runner := &PolicyRunner{
		Image:    config.Image,
		Policies: policies,
	}
	return runner, nil
}

// ExecutePolicies runs all policies stored in the policy runner.
func (pr *PolicyRunner) ExecutePolicies() {
	pr.Results.TestedImage = pr.Image
	for _, policy := range pr.Policies {
		passed, err := policy.Validate(pr.Image)

		if err != nil {
			pr.Results.Errors = append(pr.Results.Errors, policy)
			continue
		}

		if !passed {
			pr.Results.Failed = append(pr.Results.Failed, policy)
			continue
		}
		pr.Results.Passed = append(pr.Results.Passed, policy)
	}
}

// GetResults will return the results of policy execution
func (pr *PolicyRunner) GetResults() Results {
	return pr.Results
}
