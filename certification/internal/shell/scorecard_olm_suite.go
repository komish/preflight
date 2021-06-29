package shell

import (
    "os/exec"
    "os"
    "fmt"
    "io/ioutil"

    "github.com/komish/preflight/certification"
    "github.com/sirupsen/logrus"
)

type ScorecardOlmSuiteTest struct{}

func (p *ScorecardOlmSuiteTest) Validate(bundle string, logger *logrus.Logger) (bool, error) {
    
    //Curentdir
    homedir, err := os.Getwd()
    artifactsdir := homedir+"/artifacts"
    test := "operator_bundle_scorecard_OlmSuiteTest.json"
    
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
    _, err = exec.Command("oc", "image", "extract", bundle).CombinedOutput()
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
    
    scorecard_stdouterr, scorecard_err := exec.Command("operator-sdk", "scorecard",   
                                  "--config",destinationfile,
                                  "--selector=suite=olm",
                                  "--output","json","test-operator-bundle").CombinedOutput()
    
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

    err = ioutil.WriteFile(artifactsdir+"/"+test, scorecard_stdouterr, 0644)
    if err != nil {
        logger.Error("unable to copy result to /artifacts subdir: ", err)
        return false, err
    }
    if scorecard_err != nil {
        logger.Error("One or more operator-sdk scorecard test failed: ", scorecard_err)
		logger.Error("Complete scorecard log :  artifact/"+test)
		return false, scorecard_err
    }
    return true, nil
}

func (p *ScorecardOlmSuiteTest) Name() string {
    return "ScorecardOlmSuiteTest"
}

func (p *ScorecardOlmSuiteTest) Metadata() certification.Metadata {
    return certification.Metadata{
        Description:      "This test make sure that all CRs have a spec block.",
        Level:            "best",
        KnowledgeBaseURL: "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#overview", // Placeholder
        CheckURL:        "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#olm-test-suite",
    }
}

func (p *ScorecardOlmSuiteTest) Help() certification.HelpText {
    return certification.HelpText{
        Message:    "Operator-sdk scorecard OLM Test Suite. One more more test failed.",
        Suggestion: "See scorecard output artifact/operator_bundle_scorecard_OlmSuiteTest.json",
    }
}

