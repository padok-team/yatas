package volumes

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfAllVolumesHaveSnapshots(checkConfig yatas.CheckConfig, snapshot2Volumes couple, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "EC2 Volumes Snapshots"
	check.Id = testName
	check.Description = "Check if all volumes have snapshots"
	check.Status = "OK"
	for _, volume := range snapshot2Volumes.Volume {
		ok := false
		for _, snapshot := range snapshot2Volumes.Snapshot {
			if *snapshot.VolumeId == *volume.VolumeId {
				Message := "Volume " + *volume.VolumeId + " has snapshot " + *snapshot.SnapshotId
				result := results.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
				check.AddResult(result)
				ok = true
				break
			}
		}
		if !ok {
			Message := "Volume " + *volume.VolumeId + " has no snapshot"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
