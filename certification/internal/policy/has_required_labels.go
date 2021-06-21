package policy

type HasRequiredLabelPolicy struct {
}

func (p HasRequiredLabelPolicy) Validate(image string) (bool, error) {
	return true, nil
}

func (p HasRequiredLabelPolicy) GetName() string {
	return "HasRequiredLabel"
}

func (p HasRequiredLabelPolicy) GetMetadata() Metadata {
	return Metadata{
		Description:      "Checking if the container's base image is based on UBI",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		PolicyURL:        "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p HasRequiredLabelPolicy) GetHelp() HelpText {
	return HelpText{
		Message:    "It is recommened that your image be based upon the Red Hat Universal Base Image (UBI)",
		Suggestion: "Change the FROM directive in your Dockerfile or Containerfile to FROM registry.access.redhat.com/ubi8/ubi",
	}
}
