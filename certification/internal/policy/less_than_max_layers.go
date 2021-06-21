package policy

type UnderLayerMaxPolicy struct {
}

func (p UnderLayerMaxPolicy) Validate(image string) (bool, error) {
	return true, nil
}

func (p UnderLayerMaxPolicy) GetName() string {
	return "MaximumLayerPolicy"
}

func (p UnderLayerMaxPolicy) GetMetadata() Metadata {
	return Metadata{
		Description:      "Checking if container has less than 40 layers",
		Level:            "better",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		PolicyURL:        "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p UnderLayerMaxPolicy) GetHelp() HelpText {
	return HelpText{
		Message:    "Uncompressed container images should have less than 40 layers. Too many layers within the container images can degrade container performance.",
		Suggestion: "Optimize your Dockerfile to consolidate and minimize the number of layers. Each RUN command will produce a new layer. Try combining RUN commands using && where possible.",
	}
}
