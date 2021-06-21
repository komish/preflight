package policy

type RunAsNonRootPolicy struct {
}

func (p RunAsNonRootPolicy) Validate(image string) (bool, error) {
	return true, nil
}

func (p RunAsNonRootPolicy) GetName() string {
	return "RunAsNonRoot"
}

func (p RunAsNonRootPolicy) GetMetadata() Metadata {
	return Metadata{
		Description:      "Checking if container runs as the root user",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		PolicyURL:        "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p RunAsNonRootPolicy) GetHelp() HelpText {
	return HelpText{
		Message:    "A container that does not specify a non-root user will fail the automatic certification, and will be subject to a manual review before the container can be approved for publication",
		Suggestion: "Indicate a specific USER in the dockerfile or containerfile",
	}
}
