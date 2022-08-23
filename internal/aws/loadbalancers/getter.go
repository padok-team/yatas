package loadbalancers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type LoadBalancerAttributes struct {
	LoadBalancerArn  string
	LoadBalancerName string
	Output           *elasticloadbalancingv2.DescribeLoadBalancerAttributesOutput
}

func GetLoadBalancersAttributes(s aws.Config, loadbalancers []types.LoadBalancer) []LoadBalancerAttributes {
	svc := elasticloadbalancingv2.NewFromConfig(s)
	var loadBalancerAttributes []LoadBalancerAttributes
	for _, loadbalancer := range loadbalancers {
		input := &elasticloadbalancingv2.DescribeLoadBalancerAttributesInput{
			LoadBalancerArn: loadbalancer.LoadBalancerArn,
		}
		result, err := svc.DescribeLoadBalancerAttributes(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		loadBalancerAttributes = append(loadBalancerAttributes, LoadBalancerAttributes{
			LoadBalancerArn:  *loadbalancer.LoadBalancerArn,
			LoadBalancerName: *loadbalancer.LoadBalancerName,
			Output:           result,
		})
	}
	return loadBalancerAttributes
}

func GetElasticLoadBalancers(s aws.Config) []types.LoadBalancer {
	svc := elasticloadbalancingv2.NewFromConfig(s)
	var loadBalancers []types.LoadBalancer
	input := &elasticloadbalancingv2.DescribeLoadBalancersInput{
		PageSize: aws.Int32(100),
	}
	result, err := svc.DescribeLoadBalancers(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	loadBalancers = append(loadBalancers, result.LoadBalancers...)
	for {
		if result.NextMarker != nil {
			input.Marker = result.NextMarker
			result, err = svc.DescribeLoadBalancers(context.TODO(), input)
			if err != nil {
				panic(err)
			}
			loadBalancers = append(loadBalancers, result.LoadBalancers...)
		} else {
			break
		}
	}
	return loadBalancers
}
