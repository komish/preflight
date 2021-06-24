package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/komish/preflight/certification"
	"github.com/komish/preflight/certification/errors"
	"github.com/komish/preflight/certification/runtime"
)

type PolicyEngine struct {
	Image    string
	Policies []certification.Policy

	logmap         map[string][]byte
	results        runtime.Results
	localImagePath string
	isDownloaded   bool
}

// ExecutePolicies runs all policies stored in the policy engine.
func (e *PolicyEngine) ExecutePolicies() {
	e.logmap = make(map[string][]byte)

	for _, policy := range e.Policies {
		e.results.TestedImage = e.Image
		targetImage := e.Image
		execution := certification.PolicyWithRuntimeLog{
			Policy: policy,
		}

		// check if the image needs downloading
		if !e.isDownloaded {
			// TODO: consider returning an error in the ExecutePolicies instead
			// of logging those policies as errors in the output.
			isRemote, err := e.ContainerIsRemote(e.Image)
			if err != nil {
				execution.Log = []byte(err.Error())
				e.results.Errors = append(e.results.Errors, execution)
				continue
			}

			var localImagePath string
			if isRemote {
				imageTarballPath, err := e.GetContainerFromRegistry(e.Image)
				if err != nil {
					execution.Log = []byte(err.Error())
					e.results.Errors = append(e.results.Errors, execution)
					continue
				}

				localImagePath, err = e.ExtractContainerTar(imageTarballPath)
				if err != nil {
					execution.Log = []byte(err.Error())
					e.results.Errors = append(e.results.Errors, execution)
					continue
				}
			}

			e.localImagePath = localImagePath
		}

		// if we downloaded an image to disk, lets test against that.
		if len(e.localImagePath) == 0 {
			targetImage = e.localImagePath
		}

		// run the validation
		passed, logdata, err := policy.Validate(targetImage)
		e.writeToLogs(policy.Name()+".txt", logdata)
		execution.Log = logdata

		if err != nil {
			e.writeToLogs(policy.Name()+"-error.txt", []byte(err.Error()))
			e.results.Errors = append(e.results.Errors, execution)
			continue
		}

		if !passed {
			e.results.Failed = append(e.results.Failed, execution)
			continue
		}

		e.results.Passed = append(e.results.Passed, execution)
	}
}

// StorePolicy stores a given policy that needs to be executed in the policy runner.
func (e *PolicyEngine) StorePolicies(policies ...certification.Policy) {
	e.Policies = append(e.Policies, policies...)
}

// Results will return the results of policy execution.
func (e *PolicyEngine) Results() runtime.Results {
	return e.results
}

func (e *PolicyEngine) ExtractContainerTar(tarball string) (string, error) {
	// we assume the input path is something like "abcdefg.tar", representing a container image,
	// so we need to remove the extension.
	containerIDSlice := strings.Split(tarball, ".tar")
	if len(containerIDSlice) != 2 {
		// we expect a single entry in the slice, otherwise we split incorrectly
		return "", fmt.Errorf("%w: %s: %s", errors.ErrExtractingTarball, "received an improper container tarball name to extract", tarball)
	}

	outputDir := containerIDSlice[0]
	err := os.Mkdir(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrExtractingTarball, err)
	}

	stdouterr, err := exec.Command("tar", "xvf", tarball, "--directory", outputDir).CombinedOutput()
	e.writeToLogs("container-extraction.txt", stdouterr)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrExtractingTarball, err)
	}

	return outputDir, nil
}

func (e *PolicyEngine) GetContainerFromRegistry(containerLoc string) (string, error) {
	stdouterr, err := exec.Command("podman", "pull", containerLoc).CombinedOutput()
	e.writeToLogs("container-download-and-save.txt", stdouterr)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrGetRemoteContainerFailed, err)
	}
	lines := strings.Split(string(stdouterr), "\n")

	imgSig := lines[len(lines)-2]
	stdouterr, err = exec.Command("podman", "save", containerLoc, "--output", imgSig+".tar").CombinedOutput()
	e.writeToLogs("container-download-and-save.txt", stdouterr)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrSaveContainerFailed, err)
	}

	e.isDownloaded = true
	return imgSig + ".tar", nil
}

func (e *PolicyEngine) ContainerIsRemote(path string) (bool, error) {
	// TODO: Implement, for not this is just returning
	// that the resource is remote and needs to be pulled.
	return true, nil
}

func (e *PolicyEngine) Logs() map[string][]byte {
	return e.logmap
}

// writeToLogs will take the provided data and write it to a logmap to
// be written at a later time. If an empty string is provided for the
// filename, the defaultLogName will be used. Logs written to the same
// filename will be concatenated and appended.
func (e *PolicyEngine) writeToLogs(filename string, newData []byte) {

	targetLogFile := "preflight.log"

	if len(filename) != 0 {
		targetLogFile = filename
	}

	dataToWrite := newData
	existingData, exists := e.logmap[targetLogFile]
	if exists {
		dataToWrite = bytes.Join([][]byte{existingData, newData}, []byte("\n"))
	}

	e.logmap[targetLogFile] = dataToWrite
}
