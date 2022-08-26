package cloudtrail

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfCloudtrailsMultiRegion(checkConfig yatas.CheckConfig, cloudtrails []types.Trail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudtrails are in multiple regions", "check if all cloudtrails are multi region", testName)
	for _, cloudtrail := range cloudtrails {
		if !*cloudtrail.IsMultiRegionTrail {
			Message := "Cloudtrail " + *cloudtrail.Name + " is not multi region"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		} else {
			Message := "Cloudtrail " + *cloudtrail.Name + " is multi region"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudtrail.TrailARN}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
