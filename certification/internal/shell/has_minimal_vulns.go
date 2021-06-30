package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/komish/preflight/certification"
	"github.com/komish/preflight/certification/internal/utils/fileutils"
	"github.com/sirupsen/logrus"
)

const (
	ovalFilename = "rhel-8.oval.xml.bz2"
	ovalUrl      = "https://www.redhat.com/security/data/oval/v2/RHEL8/"
	reportFile   = "vuln.html"
)

type HasMinimalVulnerabilitiesCheck struct{}

func (p *HasMinimalVulnerabilitiesCheck) Validate(image string, logger *logrus.Logger) (bool, error) {

	ovalFileUrl := fmt.Sprintf("%s%s", ovalUrl, ovalFilename)

	tempdir, err := os.MkdirTemp(os.TempDir(), "oval-*")
	if err != nil {
		logger.Error("Unable to create temp dir", err)
		return false, err
	}
	logger.Debugf("Oval file dir: %s", tempdir)
	defer os.RemoveAll(filepath.Dir(tempdir))

	// TODO claen up the names
	ovalFilePath := filepath.Join(tempdir, ovalFilename)
	logger.Debugf("Oval file path: %s", ovalFilePath)

	err = fileutils.DownloadFile(ovalFilePath, ovalFileUrl)
	if err != nil {
		logger.Error("Unable to download Oval file", err)
		return false, err
	}
	// get the file name
	r := regexp.MustCompile(`(?P<filename>.*).bz2`)
	ovalFilePathDecompressed := filepath.Join(tempdir, r.FindStringSubmatch(ovalFilename)[1])

	err = fileutils.Unzip(ovalFilePath, ovalFilePathDecompressed)
	if err != nil {
		logger.Error("Unable to unzip Oval file: ", err)
		return false, err
	}

	numOfVulns, err := numberOfVulnerabilities(image, ovalFilePathDecompressed, logger)
	if err != nil {
		return false, err
	}

	logger.Debugf("The number of found vulnerabilities: %d", numOfVulns)
	if numOfVulns > 0 {
		return false, nil
	}

	return true, nil
}

func numberOfVulnerabilities(image string, ovalFilePathDecompressed string, logger *logrus.Logger) (int, error) {

	cmd := exec.Command("sudo", "oscap-podman", image, "oval", "eval", "--report", reportFile, ovalFilePathDecompressed)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		logger.Error("unable to execute oscap-podman on the image: ", cmd.Stderr)
		return 0, err
	}
	//lines := strings.Split(string(out.String()), "\n")
	r := regexp.MustCompile("Definition oval:com.redhat.*: true")
	matches := r.FindAllStringIndex(string(out.String()), -1)
	return len(matches), nil
}

func (p *HasMinimalVulnerabilitiesCheck) Name() string {
	return "HasMinimalVulnerabilities"
}

func (p *HasMinimalVulnerabilitiesCheck) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Checking for critical or important security vulnerabilites.",
		Level:            "good",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		CheckURL:         "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *HasMinimalVulnerabilitiesCheck) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "Components in the container image cannot contain any critical or important vulnerabilities, as defined at https://access.redhat.com/security/updates/classification",
		Suggestion: "Update your UBI image to the latest version or update the packages in your image to the latest versions distrubuted by Red Hat.",
	}
}
