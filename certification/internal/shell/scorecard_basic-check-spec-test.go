package shell

import (
    "os/exec"
    "os"
    "fmt"
    "io/ioutil"

    "github.com/komish/preflight/certification"
    "github.com/sirupsen/logrus"
)

type ScorecardBasicCheckSpecTest struct{}

func (p *ScorecardBasicCheckSpecTest) Validate(bundle string, logger *logrus.Logger) (bool, error) {
    
    //Curentdir
    homedir, err := os.Getwd()
    artifactsdir := homedir+"/artifacts"
    test := "operator_bundle_scorecard_BasicCheckSpecTest.json"
    
    if err != nil {
        logger.Error("unable to get current directory: ", err)
        return false, err
    }

    tmpdir, err := os.MkdirTemp("", "preflight")

    if err != nil {
        logger.Error("unable to create TempDir: ", err)
        return false, err
    }
    defer os.RemoveAll(tmpdir)

    operatordir := tmpdir+"/test-operator-bundle"
    bundledir := operatordir+"/bundle"
    
    err = os.MkdirAll(bundledir, 0777)
    if err != nil {
        logger.Error("unable to create subdirs: ", err)
        return false, err
    }

    err = os.Chdir(bundledir)
    if err != nil {
        logger.Error("unable to change dir: ", err)
        return false, err
    }

    stdouterr, err := exec.Command("oc", "image", "extract", bundle).CombinedOutput()
    if err != nil {
        logger.Error("unable to execute oc image extract on the image: ", err)
        return false, err
    }

    err = os.Chdir(tmpdir)
    if err != nil {
        logger.Error("unable to change dir: ", err)
        return false, err
    }
    
    configfile := bundledir+"/tests/scorecard/config.yaml"

    //copy config file to the root
    input, err := ioutil.ReadFile(configfile)
    if err != nil {
        logger.Error("unable to read config.yaml: ", err)
        return false, err
    }
    
    destinationfile := tmpdir+"/config.yaml"
    err = ioutil.WriteFile(destinationfile, input, 0777)
    if err != nil {
        logger.Error("unable to write to config.yaml file: ", err)
        return false, err
    }
    
    _, err = exec.Command("chmod","-R","go+r","./").CombinedOutput()
    if err != nil {
        logger.Error("unable to execute chmod: ", err)
        return false, err
    }

    fmt.Println("Running Scorecard test for ",bundle)
    
    stdouterr, err = exec.Command("operator-sdk", "scorecard",   
                                  "--config",destinationfile,
                                  "--selector=test=basic-check-spec-test",
                                  "--output","json","test-operator-bundle").CombinedOutput()
    if err != nil {
        logger.Error("unable to execute operator-sdk scorecard: ", err)
        return false, err
    }

    err = os.Chdir(homedir)
    if err != nil {
        logger.Error("unable to change dir: ", err)
        return false, err
    }

    err = os.MkdirAll(artifactsdir, 0777)
    if err != nil {
        logger.Error("unable to create artifactsdir: ", err)
        return false, err
    }

    err = ioutil.WriteFile(artifactsdir+"/"+test, stdouterr, 0644)
    if err != nil {
        logger.Error("unable to copy result to /artifacts subdir: ", err)
        return false, err
    }
    
    return true, nil
}

func (p *ScorecardBasicCheckSpecTest) Name() string {
    return "ScorecardBasicCheckSpecTest"
}

func (p *ScorecardBasicCheckSpecTest) Metadata() certification.Metadata {
    return certification.Metadata{
        Description:      "This test make sure that all CRs have a spec block.",
        Level:            "best",
        KnowledgeBaseURL: "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#overview", // Placeholder
        CheckURL:        "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#basic-test-suite",
    }
}

func (p *ScorecardBasicCheckSpecTest) Help() certification.HelpText {
    return certification.HelpText{
        Message:    "Operator-sdk scorecard basic-check-spec-test failed.",
        Suggestion: "Make sure that all CRs have a spec block",
    }
}

