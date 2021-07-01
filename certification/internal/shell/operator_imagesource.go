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
		logger.Error(" Failed to execute cmdchain builder")
		logger.Debug(" failed to execute cmdchain builder", err)
		return false, nil
	}

	//s := b.String() //cast to string
	userRegistry := strings.TrimRight(output.String(), "\n")
	logger.Info("Check Image registry for : ", userRegistry)
	
	if approvedRegistries[userRegistry] {
		logger.Debug(userRegistry, " found in approved registry")
		return true, nil
	}

	logger.Info(userRegistry," not found in approved registry")
	return false, nil
}

func (p *OperatorImageSourceCheck) Name() string {
	return "OperatorImageSourceCheck"
}

func (p *OperatorImageSourceCheck) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Validating Imagesource",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		CheckURL:         "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *OperatorImageSourceCheck) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "Imagesource check failed! Non-approved images found.",
		Suggestion: "Push image to one of the approved registries.\n",
	}
}
