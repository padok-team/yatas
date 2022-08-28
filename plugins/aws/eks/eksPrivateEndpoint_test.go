package eks

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfEksEndpointPrivate(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		clusters    []types.Cluster
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name:     aws.String("test"),
						Endpoint: aws.String("https://test.eks.amazonaws.com"),
						ResourcesVpcConfig: &types.VpcConfigResponse{
							EndpointPrivateAccess: true,
						},
					},
				},
			},
			want: "OK",
		},
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name:     aws.String("test"),
						Endpoint: aws.String("https://test.eks.amazonaws.com"),
						ResourcesVpcConfig: &types.VpcConfigResponse{
							EndpointPrivateAccess: true,
							EndpointPublicAccess:  true,
							PublicAccessCidrs:     []string{"0.0.0.0/0"},
						},
					},
				},
			},
			want: "FAIL",
		},
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name:     aws.String("test"),
						Endpoint: aws.String("https://test.eks.amazonaws.com"),
						ResourcesVpcConfig: &types.VpcConfigResponse{
							EndpointPrivateAccess: true,
							EndpointPublicAccess:  true,
							PublicAccessCidrs:     []string{"0.0.0.0/8"},
						},
					},
				},
			},
			want: "OK",
		},
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name:     aws.String("test"),
						Endpoint: aws.String("https://test.eks.amazonaws.com"),
					},
				},
			},
			want: "FAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfEksEndpointPrivate(tt.args.checkConfig, tt.args.clusters, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfEksEndpointPrivate() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
