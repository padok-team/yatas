package volumes

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetSnapshots(s *session.Session) []*ec2.Snapshot {
	svc := ec2.New(s)
	input := &ec2.DescribeSnapshotsInput{}
	result, err := svc.DescribeSnapshots(input)
	if err != nil {
		panic(err)
	}
	return result.Snapshots
}

func CheckIfAllVolumesHaveSnapshots(s *session.Session, volumes []*ec2.Volume, c *[]types.Check) {
	logger.Info("Running AWS_VOL_002")
	var check types.Check
	check.Name = "EC2 Volumes Snapshots"
	check.Id = "AWS_VOL_002"
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
