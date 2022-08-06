package volumes

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetSnapshots(s *session.Session) []*ec2.Snapshot {
	svc := ec2.New(s)
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: []*string{aws.String("self")},
	}
	result, err := svc.DescribeSnapshots(input)
	if err != nil {
		panic(err)
	}
	return result.Snapshots
}

func CheckIfAllVolumesHaveSnapshots(s *session.Session, volumes []*ec2.Volume, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Volumes Snapshots"
	check.Id = testName
	check.Description = "Check if all volumes have snapshots"
	check.Status = "OK"
	snapshots := GetSnapshots(s)
	for _, volume := range volumes {
		ok := false
		for _, snapshot := range snapshots {
			if *snapshot.VolumeId == *volume.VolumeId {
				status := "OK"
				Message := "Volume " + *volume.VolumeId + " has snapshot " + *snapshot.SnapshotId
				check.Results = append(check.Results, types.Result{Status: status, Message: Message})
				ok = true
				break
			}
		}
		if !ok {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Volume " + *volume.VolumeId + " has no snapshot"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfAllSnapshotsEncrypted(s *session.Session, snapshots []*ec2.Snapshot, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Snapshots Encryption"
	check.Id = testName
	check.Description = "Check if all snapshots are encrypted"
	check.Status = "OK"
	for _, snapshot := range snapshots {
		if *snapshot.Encrypted == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Snapshot " + *snapshot.SnapshotId + " is not encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Snapshot " + *snapshot.SnapshotId + " is encrypted"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfSnapshotYoungerthan24h(s *session.Session, snapshots []*ec2.Snapshot, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Snapshots Age"
	check.Id = testName
	check.Description = "Check if all snapshots are younger than 24h"
	check.Status = "OK"
	for _, snapshot := range snapshots {
		creationTime := *snapshot.StartTime
		if creationTime.After(time.Now().Add(-24 * time.Hour)) {
			status := "OK"
			Message := "Snapshot " + *snapshot.SnapshotId + " is younger than 24h"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Snapshot " + *snapshot.SnapshotId + " is older than 24h"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}
