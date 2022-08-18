package volumes

import (
	"fmt"
	"time"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfSnapshotYoungerthan24h(checkConfig yatas.CheckConfig, vs couple, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Snapshots Age", "Check if all snapshots are younger than 24h", testName)
	for _, volume := range vs.Volume {
		snapshotYoungerThan24h := false
		for _, snapshot := range vs.Snapshot {
			if *snapshot.VolumeId == *volume.VolumeId {
				creationTime := *snapshot.StartTime
				if creationTime.After(time.Now().Add(-24 * time.Hour)) {
					snapshotYoungerThan24h = true
					break
				}
			}
		}
		if !snapshotYoungerThan24h {
			Message := "Volume " + *volume.VolumeId + " has no snapshot younger than 24h"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.Results = append(check.Results, result)
		} else {
			Message := "Volume " + *volume.VolumeId + " has snapshot younger than 24h"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
