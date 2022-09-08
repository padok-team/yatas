package volumes

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type couple struct {
	Volume   []types.Volume
	Snapshot []types.Snapshot
}

// GetSnapshots returns all snapshots for an aws config
func GetSnapshots(s aws.Config) []types.Snapshot {
	svc := ec2.NewFromConfig(s)
	var snapshots []types.Snapshot
	input := &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{*aws.String("self")},
	}
	result, err := svc.DescribeSnapshots(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	snapshots = append(snapshots, result.Snapshots...)
	for {
		if result.NextToken != nil {
			input.NextToken = result.NextToken
			result, err = svc.DescribeSnapshots(context.TODO(), input)
			snapshots = append(snapshots, result.Snapshots...)
			if err != nil {
				panic(err)
			}
		} else {
			break
		}
	}
	return snapshots
}

// GetVolumesAndSnapshots returns all volumes  for an aws config
func GetVolumes(s aws.Config) []types.Volume {
	svc := ec2.NewFromConfig(s)
	var volumes []types.Volume
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	volumes = append(volumes, result.Volumes...)
	for {
		if result.NextToken != nil {
			input.NextToken = result.NextToken
			result, err = svc.DescribeVolumes(context.TODO(), input)
			volumes = append(volumes, result.Volumes...)
			if err != nil {
				panic(err)
			}
		} else {
			break
		}

	}
	return volumes
}
