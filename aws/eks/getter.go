package eks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EKSGetObjectAPI interface {
	ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
}

func GetClusters(svc EKSGetObjectAPI) []types.Cluster {
	input := &eks.ListClustersInput{}
	result, err := svc.ListClusters(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var clusters []string
	var clustersDetails []types.Cluster
	for _, r := range result.Clusters {
		clusters = append(clusters, r)
	}
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.ListClusters(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		for _, r := range result.Clusters {
			clusters = append(clusters, r)
		}
	}

	for _, c := range clusters {
		input := &eks.DescribeClusterInput{
			Name: aws.String(c),
		}
		result, err := svc.DescribeCluster(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		clustersDetails = append(clustersDetails, *result.Cluster)
	}
	return clustersDetails

}
