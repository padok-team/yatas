package volumes

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetVolumes(s *session.Session) []*ec2.Volume {
	svc := ec2.New(s)
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(input)
	if err != nil {
		panic(err)
	}
	return result.Volumes
}

func checkIfEncryptionEnabled(s *session.Session, volumes []*ec2.Volume, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Volumes Encryption"
	check.Id = testName
	check.Description = "Check if EC2 encryption is enabled"
	check.Status = "OK"
	svc := ec2.New(s)
	for _, volume := range volumes {
		params := &ec2.DescribeVolumesInput{
			VolumeIds: []*string{volume.VolumeId},
		}
		resp, err := svc.DescribeVolumes(params)
		if err != nil {
			panic(err)
		}
		if *resp.Volumes[0].Encrypted == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func CheckIfVolumesTypeGP3(s *session.Session, volumes []*ec2.Volume, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "EC2 Volumes Type"
	check.Id = testName
	check.Description = "Check if all volumes are of type gp3"
	check.Status = "OK"
	for _, volume := range volumes {
		if *volume.VolumeType != "gp3" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

type couple struct {
	volume   []*ec2.Volume
	snapshot []*ec2.Snapshot
}

func RunVolumesTest(s *session.Session, c *yatas.Config) []types.Check {
	var checks []types.Check
	logger.Debug("Starting EC2 volumes tests")
	volumes := GetVolumes(s)
	snapshots := GetSnapshots(s)
	couples := couple{volumes, snapshots}

	yatas.CheckTest(c, "AWS_VOL_001", checkIfEncryptionEnabled)(s, volumes, "AWS_VOL_001", &checks)
	yatas.CheckTest(c, "AWS_VOL_002", CheckIfVolumesTypeGP3)(s, volumes, "AWS_VOL_002", &checks)
	yatas.CheckTest(c, "AWS_VOL_003", CheckIfAllVolumesHaveSnapshots)(s, volumes, "AWS_VOL_004", &checks)

	yatas.CheckTest(c, "AWS_BAK_001", CheckIfAllSnapshotsEncrypted)(s, snapshots, "AWS_BAK_001", &checks)
	yatas.CheckTest(c, "AWS_BAK_002", CheckIfSnapshotYoungerthan24h)(s, couples, "AWS_BAK_002", &checks)

	return checks
}
