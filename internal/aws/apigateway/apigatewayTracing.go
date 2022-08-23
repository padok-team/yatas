package apigateway

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfTracingEnabled(checkConfig yatas.CheckConfig, stages []types.Stage, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Apigateway Tracing enabled", "Check if all stages are enabled for tracing", testName)
	for _, stage := range stages {
		if stage.TracingEnabled {
			Message := "Tracing is enabled on stage" + *stage.StageName
			result := results.Result{Status: "OK", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		} else {
			Message := "Tracing is not enabled on " + *stage.StageName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *stage.StageName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
