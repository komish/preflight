package shell

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/itchyny/gojq"
	"github.com/komish/preflight/certification"
	"github.com/sirupsen/logrus"
)

var maxLayers int = 40

type UnderLayerMaxCheck struct {
}

func (p *UnderLayerMaxCheck) Validate(image string, logger *logrus.Logger) (bool, error) {
	stdouterr, err := exec.Command("podman", "inspect", image).CombinedOutput()
	if err != nil {
		logger.Error("unable to execute inspect on the image: ", err)
		return false, err
	}
	var inspectData []interface{}
	err = json.Unmarshal(stdouterr, &inspectData)
	if err != nil {
		logger.Error("unable to parse podman inspect data for image", err)
		logger.Debug("error marshaling podman inspect data: ", err)
		logger.Trace("failure in attempt to convert the raw bytes from `podman inspect` to a []interface{}")
		return false, err
	}

	jqQueryString := ".[0].RootFS.Layers"

	query, err := gojq.Parse(jqQueryString)
	if err != nil {
		logger.Error("unable to parse podman inspect data for image", err)
		logger.Debug("unable to successfully parse the gojq query string:", err)
		return false, err
	}

	iter := query.Run(inspectData)
	val, nextOk := iter.Next()

	if !nextOk {
		logger.Warn("did not receive any layer information when parsing container image")
		return false, nil
	}
	if err, ok := val.(error); ok {
		logger.Error("unable to parse podman inspect data for image", err)
		logger.Debug("unable to successfully parse the podman inspect output with the query string provided:", err)
		// this is an error, as we didn't get the proper input from `podman inspect`
		return false, err
	}

	layers := val.([]interface{})

	if len(layers) > maxLayers {
		return false, nil
	}
	return true, nil
}

func (p *UnderLayerMaxCheck) Name() string {
	return "UnderMaxLayers"
}

func (p *UnderLayerMaxCheck) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      fmt.Sprintf("Checking if container has less than %d layers", maxLayers),
		Level:            "better",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		CheckURL:         "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *UnderLayerMaxCheck) Help() certification.HelpText {
	return certification.HelpText{
		Message:    fmt.Sprintf("Uncompressed container images should have less than %d layers. Too many layers within the container images can degrade container performance.", maxLayers),
		Suggestion: "Optimize your Dockerfile to consolidate and minimize the number of layers. Each RUN command will produce a new layer. Try combining RUN commands using && where possible.",
	}
}
