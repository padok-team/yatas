package volumes

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfVolumeIsUsed(checkConfig yatas.CheckConfig, volumes []types.Volume, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("EC2's volumes are unused", "Check if EC2 volumes are unused", testName)
	for _, volume := range volumes {
		if volume.State != types.VolumeStateInUse && volume.State != types.VolumeStateDeleted {
			Message := "EC2 volume is unused " + *volume.VolumeId
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else if volume.State == types.VolumeStateDeleted {
			continue
		} else {
			Message := "EC2 volume is in use " + *volume.VolumeId
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
