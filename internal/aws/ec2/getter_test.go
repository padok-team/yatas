package ec2

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type MockSVCEC2 func()

func (m MockSVCEC2) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	token := aws.String("ididididid")
	if params.NextToken != nil {
		token = nil
	}

	return &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId:   aws.String("instanceId"),
						InstanceType: types.InstanceTypeA12xlarge,
					},
				},
			},
		},
		NextToken: token,
	}, nil
}

func TestGetEC2s(t *testing.T) {
	type args struct {
		svc EC2GetObjectAPI
	}
	tests := []struct {
		name string
		args args
		want []types.Instance
	}{
		{
			name: "TestGetEC2s",
			args: args{
				svc: MockSVCEC2(nil),
			},
			want: []types.Instance{
				{
					InstanceId:   aws.String("instanceId"),
					InstanceType: types.InstanceTypeA12xlarge,
				},
				{
					InstanceId:   aws.String("instanceId"),
					InstanceType: types.InstanceTypeA12xlarge,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEC2s(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEC2s() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
