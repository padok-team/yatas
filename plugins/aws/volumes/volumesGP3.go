package volumes

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfVolumesTypeGP3(checkConfig yatas.CheckConfig, volumes []types.Volume, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("EC2 are using GP3", "Check if all volumes are of type gp3", testName)
	for _, volume := range volumes {
		if volume.VolumeType != "gp3" {
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
