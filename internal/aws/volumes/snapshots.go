package volumes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
)

func GetSnapshots(s aws.Config) []types.Snapshot {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{*aws.String("self")},
	}
	result, err := svc.DescribeSnapshots(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Snapshots
}

func CheckIfAllVolumesHaveSnapshots(wg *sync.WaitGroup, s aws.Config, volumes []types.Volume, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "EC2 Volumes Snapshots"
	check.Id = testName
	check.Description = "Check if all volumes have snapshots"
	check.Status = "OK"
	snapshots := GetSnapshots(s)
	for _, volume := range volumes {
		ok := false
		for _, snapshot := range snapshots {
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
	queueToAdd <- check
}

func CheckIfAllSnapshotsEncrypted(wg *sync.WaitGroup, s aws.Config, snapshots []types.Snapshot, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Snapshots Encryption", "Check if all snapshots are encrypted", testName)
	for _, snapshot := range snapshots {
		if !*snapshot.Encrypted {
			Message := "Snapshot " + *snapshot.SnapshotId + " is not encrypted"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *snapshot.SnapshotId}
			check.AddResult(result)
		} else {
			Message := "Snapshot " + *snapshot.SnapshotId + " is encrypted"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *snapshot.SnapshotId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfSnapshotYoungerthan24h(wg *sync.WaitGroup, s aws.Config, vs couple, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Snapshots Age", "Check if all snapshots are younger than 24h", testName)
	for _, volume := range vs.volume {
		snapshotYoungerThan24h := false
		for _, snapshot := range vs.snapshot {
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
	queueToAdd <- check
}
