package shell

import (
    "os/exec"
    "os"
    "strings"
    "io/ioutil"
    "path/filepath"
    "encoding/json"

    "github.com/komish/preflight/certification"
    "github.com/sirupsen/logrus"
    "github.com/itchyny/gojq"
)

type ScorecardBasicSpecCheck struct{}
const scorecardBasicCheckResult string = "operator_bundle_scorecard_BasicSpecCheck.json"

func (p *ScorecardBasicSpecCheck) Validate(bundleImage string, logger *logrus.Logger) (bool, error) {
    
    currentDir, err := os.Getwd()
    if err != nil {
        logger.Error("unable to get current directory: ", err)
        return false, err
    }

    artifactsDir := filepath.Join(currentDir, "/artifacts")
    
    err = os.MkdirAll(artifactsDir, 0777)
    if err != nil {
        logger.Error("unable to create artifactsDir: ", err)
        return false, err
    }
    
    logger.Debug("Running operator-sdk scorecard check for ",bundleImage)
    logger.Debug("--selector=test=basic-check-spec-test")
    stdouterr, err := exec.Command("operator-sdk", "scorecard",   
                                   "--selector=test=basic-check-spec-test",
                                   "--output", "json", bundleImage).CombinedOutput()
    
    scorecardFile := filepath.Join(artifactsDir, "/", scorecardBasicCheckResult)
    
    err = ioutil.WriteFile(scorecardFile, stdouterr, 0644)
    if err != nil {
        logger.Error("unable to copy result to /artifacts subdir: ", err)
        return false, err
    }

    // we must send gojq a interface{}, so we have to convert our inspect output to that type
    var inspectData interface{}
    
    err = json.Unmarshal(stdouterr, &inspectData)
    if err != nil {
        logger.Error("unable to parse scorecard json output")
        logger.Debug("error Unmarshaling scorecard json output: ", err)
        logger.Debug("operator_sdk failed to execute.")
        logger.Trace("failure in attempt to convert the raw bytes from `operator-sdk scorecard` to a interface{}")
        return false, err
    }

    query, err := gojq.Parse(".items[].status.results[] | .name, .state")
    if err != nil {
        logger.Error("unable to parse scorecard json output")
        logger.Debug("unable to parse :", err)
        return false, err
    }

    // gojq expects us to iterate in the event that our query returned multiple matching values, but we only expect one.
    iter := query.Run(inspectData)
    
    foundTestFailed := false

    logger.Info("scorecard output")
    
    for {
        v, ok := iter.Next()
        if !ok {
            logger.Warn("Did not receive any test result information when parsing scorecard output.")
	    // in this case, there was no data returned from jq, so we need to fail the check.
	    break
        }
        if err, ok := v.(error); ok {
            logger.Error("unable to parse scorecard output")
            logger.Debug("unable to successfully parse the scorecard output", err)
            return false, err
        }
        //test fails but keeps going listing out all tests
        s := v.(string)
        logger.Info(s)
        if strings.Contains(s, "fail") {
            foundTestFailed = true
        }
    }
    return !foundTestFailed, nil
}

func (p *ScorecardBasicSpecCheck) Name() string {
    return "ScorecardBasicSpecCheck"
}

func (p *ScorecardBasicSpecCheck) Metadata() certification.Metadata {
    return certification.Metadata{
        Description:      "Check to make sure that all CRs have a spec block.",
        Level:            "best",
        KnowledgeBaseURL: "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#overview", // Placeholder
        CheckURL:        "https://sdk.operatorframework.io/docs/advanced-topics/scorecard/scorecard/#basic-test-suite",
    }
}

func (p *ScorecardBasicSpecCheck) Help() certification.HelpText {
    return certification.HelpText{
        Message:    "Operator-sdk scorecard basic spec check failed.",
        Suggestion: "Make sure that all CRs have a spec block",
    }
}