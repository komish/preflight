package shell

import (
	"bytes"
	"strings"
	"github.com/komish/preflight/certification"
	"github.com/sirupsen/logrus"
	"github.com/rainu/go-command-chain"
)

type OperatorImageSourceCheck struct {
}

func (p *OperatorImageSourceCheck) Validate(bundleImage string, logger *logrus.Logger) (bool, error) {
	
	approvedRegistries := map[string]bool{
		"registry.connect.dev.redhat.com": true,
		"registry.connect.qa.redhat.com": true,
		"registry.connect.stage.redhat.com": true,
		"registry.connect.redhat.com": true,
		"registry.redhat.io": true,
		"registry.access.redhat.com": true,
		"quay.io": true,
	}
	
	output := &bytes.Buffer{}
    inputContent := strings.NewReader(bundleImage)
	
	err := cmdchain.Builder().
		WithInput(inputContent).
		Join("cut", "-d",",", "-f1").
		Join("cut", "-d","/", "-f1").
		Finalize().WithOutput(output).Run()
	if err != nil {
		panic(err)
	}

	//s := b.String() //cast to string
	s := strings.TrimRight(output.String(), "\n")
	logger.Info("Check Image registry for : ", s)
	
	if approvedRegistries[s] {
		logger.Debug(s, " found in the approved registry")
		return true, nil
	}
	
	return false, nil
}

func (p *OperatorImageSourceCheck) Name() string {
	return "OperatorImageSourceCheck"
}

func (p *OperatorImageSourceCheck) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Validating Bundle image",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		CheckURL:         "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *OperatorImageSourceCheck) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "Operator sdk validation test failed, this test checks if it can validate the content and format of the operator bundle",
		Suggestion: "Valid bundles are definied by bundle spec, so make sure that this bundle conforms to that spec. More Information: https://github.com/operator-framework/operator-registry/blob/master/docs/design/operator-bundle.md",
	}
}
