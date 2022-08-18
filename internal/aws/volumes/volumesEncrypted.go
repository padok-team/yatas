package volumes

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, volumes []types.Volume, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Volumes Encryption", "Check if EC2 encryption is enabled", testName)
	for _, volume := range volumes {
		if volume.Encrypted != nil && *volume.Encrypted {
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
