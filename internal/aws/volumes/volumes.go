package volumes

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetVolumes(s aws.Config) []types.Volume {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Volumes
}

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, volumes []types.Volume, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Volumes Encryption", "Check if EC2 encryption is enabled", testName)
	svc := ec2.NewFromConfig(checkConfig.ConfigAWS)
	for _, volume := range volumes {
		params := &ec2.DescribeVolumesInput{
			VolumeIds: []string{*volume.VolumeId},
		}
		resp, err := svc.DescribeVolumes(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if *resp.Volumes[0].Encrypted {
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfVolumesTypeGP3(checkConfig yatas.CheckConfig, volumes []types.Volume, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Volumes Type", "Check if all volumes are of type gp3", testName)
	for _, volume := range volumes {
		if volume.VolumeType != "gp3" {
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		} else {
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *volume.VolumeId}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

type couple struct {
	volume   []types.Volume
	snapshot []types.Snapshot
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	logger.Debug("Starting EC2 volumes tests")
	volumes := GetVolumes(s)
	snapshots := GetSnapshots(s)
	couples := couple{volumes, snapshots}

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VOL_001", checkIfEncryptionEnabled)(checkConfig, volumes, "AWS_VOL_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VOL_002", CheckIfVolumesTypeGP3)(checkConfig, volumes, "AWS_VOL_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VOL_003", CheckIfAllVolumesHaveSnapshots)(checkConfig, volumes, "AWS_VOL_003")

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_BAK_001", CheckIfAllSnapshotsEncrypted)(checkConfig, snapshots, "AWS_BAK_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_BAK_002", CheckIfSnapshotYoungerthan24h)(checkConfig, couples, "AWS_BAK_002")

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
