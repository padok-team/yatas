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

func GetVolumes(s aws.Config) []types.Volume {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeVolumesInput{}
	result, err := svc.DescribeVolumes(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Volumes
}
