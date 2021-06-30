package shell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/komish/preflight/certification/errors"
	"github.com/sirupsen/logrus"
)

type Image struct {
	Config Config
}

type Config struct {
	Labels map[string]string
}

func GetLabelsForImage(image string, logger *logrus.Logger) (*map[string]string, error) {
	cmd := exec.Command("podman", "image", "inspect", image)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrGetRemoteContainerFailed, err)
	}

	var images []Image
	err = json.Unmarshal(stdout.Bytes(), &images)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrImageInspectFailed, err)
	}

	return &images[0].Config.Labels, nil
}
