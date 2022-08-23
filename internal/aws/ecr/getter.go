package ecr

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func GetECRs(s aws.Config) []types.Repository {
	svc := ecr.NewFromConfig(s)
	var ecrRepositories []types.Repository
	input := &ecr.DescribeRepositoriesInput{
		MaxResults: aws.Int32(100),
	}
	result, err := svc.DescribeRepositories(context.TODO(), input)
	ecrRepositories = append(ecrRepositories, result.Repositories...)
	if err != nil {
		panic(err)
	}
	for {
		if result.NextToken != nil {
			input.NextToken = result.NextToken
			result, err = svc.DescribeRepositories(context.TODO(), input)
			ecrRepositories = append(ecrRepositories, result.Repositories...)
			if err != nil {
				panic(err)
			}
		} else {
			break
		}
	}

	return ecrRepositories
}
