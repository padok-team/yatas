package vpc

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetListVPC(s aws.Config) []types.Vpc {
	svc := ec2.NewFromConfig(s)
	var vpcs []types.Vpc
	input := &ec2.DescribeVpcsInput{}
	result, err := svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	vpcs = append(vpcs, result.Vpcs...)
	for {
		if result.NextToken == nil {
			break
		}
		input.NextToken = result.NextToken
		result, err = svc.DescribeVpcs(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		vpcs = append(vpcs, result.Vpcs...)
	}
	return vpcs
}

type VPCToSubnet struct {
	VpcID   string
	Subnets []types.Subnet
}

func GetSubnetForVPCS(s aws.Config, vpcs []types.Vpc) []VPCToSubnet {
	svc := ec2.NewFromConfig(s)
	var vpcSubnets []VPCToSubnet
	for _, vpc := range vpcs {
		input := &ec2.DescribeSubnetsInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("vpc-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeSubnets(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		vpcSubnets = append(vpcSubnets, VPCToSubnet{
			VpcID:   *vpc.VpcId,
			Subnets: result.Subnets,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeSubnets(context.TODO(), input)
			if err != nil {
				panic(err)
			}
			vpcSubnets = append(vpcSubnets, VPCToSubnet{
				VpcID:   *vpc.VpcId,
				Subnets: result.Subnets,
			})
		}
	}
	return vpcSubnets
}

type VpcToInternetGateway struct {
	VpcID            string
	InternetGateways []types.InternetGateway
}

func GetInternetGatewaysForVpc(s aws.Config, vpcs []types.Vpc) []VpcToInternetGateway {
	svc := ec2.NewFromConfig(s)
	var vpcInternetGateways []VpcToInternetGateway
	for _, vpc := range vpcs {
		input := &ec2.DescribeInternetGatewaysInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("attachment.vpc-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeInternetGateways(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		vpcInternetGateways = append(vpcInternetGateways, VpcToInternetGateway{
			VpcID:            *vpc.VpcId,
			InternetGateways: result.InternetGateways,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeInternetGateways(context.TODO(), input)
			if err != nil {
				panic(err)
			}
			vpcInternetGateways = append(vpcInternetGateways, VpcToInternetGateway{
				VpcID:            *vpc.VpcId,
				InternetGateways: result.InternetGateways,
			})
		}
	}
	return vpcInternetGateways
}

type VpcToFlowLogs struct {
	VpcID    string
	FlowLogs []types.FlowLog
}

func GetFlowLogsForVpc(s aws.Config, vpcs []types.Vpc) []VpcToFlowLogs {
	svc := ec2.NewFromConfig(s)
	var vpcFlowLogs []VpcToFlowLogs
	for _, vpc := range vpcs {
		input := &ec2.DescribeFlowLogsInput{
			Filter: []types.Filter{
				{
					Name:   aws.String("resource-id"),
					Values: []string{*vpc.VpcId},
				},
			},
		}
		result, err := svc.DescribeFlowLogs(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		vpcFlowLogs = append(vpcFlowLogs, VpcToFlowLogs{
			VpcID:    *vpc.VpcId,
			FlowLogs: result.FlowLogs,
		})
		for {
			if result.NextToken == nil {
				break
			}
			input.NextToken = result.NextToken
			result, err = svc.DescribeFlowLogs(context.TODO(), input)
			if err != nil {
				panic(err)
			}
			vpcFlowLogs = append(vpcFlowLogs, VpcToFlowLogs{
				VpcID:    *vpc.VpcId,
				FlowLogs: result.FlowLogs,
			})
		}
	}
	return vpcFlowLogs
}
