package eks

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type mocksvc func()

func (m mocksvc) ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	return &eks.ListClustersOutput{
		Clusters: []string{"test"},
	}, nil
}

func (m mocksvc) DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	return &eks.DescribeClusterOutput{
		Cluster: &types.Cluster{
			Name: aws.String("test"),
		},
	}, nil
}

func TestGetClusters(t *testing.T) {
	type args struct {
		svc EKSGetObjectAPI
	}
	tests := []struct {
		name string
		args args
		want []types.Cluster
	}{
		{
			name: "test",
			args: args{
				svc: mocksvc(nil),
			},
			want: []types.Cluster{
				{
					Name: aws.String("test"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetClusters(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClusters() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
