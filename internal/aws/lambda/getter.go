package lambda

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func GetLambdas(s aws.Config) []types.FunctionConfiguration {
	svc := lambda.NewFromConfig(s)
	var lambdas []types.FunctionConfiguration
	input := &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(100),
	}
	result, err := svc.ListFunctions(context.TODO(), input)
	lambdas = append(lambdas, result.Functions...)
	if err != nil {
		panic(err)
	}
	for {
		if result.NextMarker != nil {
			input.Marker = result.NextMarker
			result, err = svc.ListFunctions(context.TODO(), input)
			lambdas = append(lambdas, result.Functions...)
			if err != nil {
				panic(err)
			}
		} else {
			break
		}
	}
	return lambdas
}
