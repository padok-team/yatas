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

func checkIfEncryptionEnabled(wg *sync.WaitGroup, s aws.Config, volumes []types.Volume, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "EC2 Volumes Encryption"
	check.Id = testName
	check.Description = "Check if EC2 encryption is enabled"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, volume := range volumes {
		params := &ec2.DescribeVolumesInput{
			VolumeIds: []string{*volume.VolumeId},
		}
		resp, err := svc.DescribeVolumes(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if *resp.Volumes[0].Encrypted {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "EC2 encryption is not enabled on " + *volume.VolumeId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "EC2 encryption is enabled on " + *volume.VolumeId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfVolumesTypeGP3(wg *sync.WaitGroup, s aws.Config, volumes []types.Volume, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "EC2 Volumes Type"
	check.Id = testName
	check.Description = "Check if all volumes are of type gp3"
	check.Status = "OK"
	for _, volume := range volumes {
		if volume.VolumeType != "gp3" {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Volume " + *volume.VolumeId + " is not of type gp3"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "Volume " + *volume.VolumeId + " is of type gp3"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

type couple struct {
	volume   []types.Volume
	snapshot []types.Snapshot
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	logger.Debug("Starting EC2 volumes tests")
	volumes := GetVolumes(s)
	snapshots := GetSnapshots(s)
	couples := couple{volumes, snapshots}
	var wg sync.WaitGroup

	go yatas.CheckTest(&wg, c, "AWS_VOL_001", checkIfEncryptionEnabled)(&wg, s, volumes, "AWS_VOL_001", &checks)
	go yatas.CheckTest(&wg, c, "AWS_VOL_002", CheckIfVolumesTypeGP3)(&wg, s, volumes, "AWS_VOL_002", &checks)
	go yatas.CheckTest(&wg, c, "AWS_VOL_003", CheckIfAllVolumesHaveSnapshots)(&wg, s, volumes, "AWS_VOL_003", &checks)

	go yatas.CheckTest(&wg, c, "AWS_BAK_001", CheckIfAllSnapshotsEncrypted)(&wg, s, snapshots, "AWS_BAK_001", &checks)
	go yatas.CheckTest(&wg, c, "AWS_BAK_002", CheckIfSnapshotYoungerthan24h)(&wg, s, couples, "AWS_BAK_002", &checks)

	wg.Wait()
	if c.Progress != nil {

	}
	queue <- checks
}
