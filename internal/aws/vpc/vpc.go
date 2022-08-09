package vpc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetListVPC(s aws.Config) []types.Vpc {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeVpcsInput{}
	result, err := svc.DescribeVpcs(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Vpcs
}

func checkCIDR20(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "VPC CIDR"
	check.Id = testName
	check.Description = "Check if VPC CIDR is /20 or bigger"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeVpcsInput{
			VpcIds: []string{*vpc.VpcId},
		}
		resp, err := svc.DescribeVpcs(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		cidr := *resp.Vpcs[0].CidrBlock
		// split the cidr to / and get the last part as an int
		cidrInt, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
		if cidrInt > 20 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC CIDR is not /20 or bigger on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			status := "OK"
			Message := "VPC CIDR is /20 or bigger on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}
	*c = append(*c, check)
}

func checkIfVPCFLowLogsEnabled(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "VPC Flow Logs"
	check.Id = testName
	check.Description = "Check if VPC Flow Logs are enabled"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeFlowLogsInput{
			Filter: []types.Filter{
				{
					Name: aws.String("resource-id"),
					Values: []string{
						*vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeFlowLogs(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if len(resp.FlowLogs) == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC Flow Logs are not enabled on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			status := "OK"
			Message := "VPC Flow Logs are enabled on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}
	*c = append(*c, check)
}

func checkIfOnlyOneGateway(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "VPC Gateway"
	check.Id = testName
	check.Description = "Check if VPC has only one gateway"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeInternetGatewaysInput{
			Filters: []types.Filter{
				{
					Name: aws.String("attachment.vpc-id"),
					Values: []string{
						*vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeInternetGateways(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if len(resp.InternetGateways) > 1 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC has more than one gateway on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			status := "OK"
			Message := "VPC has only one gateway on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}
	*c = append(*c, check)
}

func checkIfOnlyOneVPC(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "VPC Only One"
	check.Id = testName
	check.Description = "Check if VPC has only one VPC"
	check.Status = "OK"
	for _, vpc := range vpcs {
		if len(vpcs) > 1 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC Id:" + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			status := "OK"
			Message := "VPC Id:" + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}

	*c = append(*c, check)
}

func CheckIfSubnetInDifferentZone(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Subnets in different zone"
	check.Id = testName
	check.Description = "Check if Subnet are in different zone"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeSubnetsInput{
			Filters: []types.Filter{
				{
					Name: aws.String("vpc-id"),
					Values: []string{
						*vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeSubnets(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		subnetsAZ := make(map[string]int)
		for _, subnet := range resp.Subnets {
			subnetsAZ[*subnet.AvailabilityZone]++
		}
		if len(subnetsAZ) > 1 {
			check.Status = "OK"
			status := "OK"
			Message := "Subnets are in different zone on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "Subnets are in same zone on " + *vpc.VpcId
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}
	*c = append(*c, check)
}

func CheckIfAtLeast2Subnets(s aws.Config, vpcs []types.Vpc, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "At least 2 subnets"
	check.Id = testName
	check.Description = "Check if VPC has at least 2 subnets"
	check.Status = "OK"
	svc := ec2.NewFromConfig(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeSubnetsInput{
			Filters: []types.Filter{
				{
					Name: aws.String("vpc-id"),
					Values: []string{
						*vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeSubnets(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if len(resp.Subnets) < 2 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC " + *vpc.VpcId + " has less than 2 subnets"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		} else {
			status := "OK"
			Message := "VPC " + *vpc.VpcId + " has at least 2 subnets"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *vpc.VpcId})
		}
	}
	*c = append(*c, check)
}

func RunVPCTests(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	vpcs := GetListVPC(s)
	yatas.CheckTest(c, "AWS_VPC_001", checkCIDR20)(s, vpcs, "AWS_VPC_001", &checks)
	yatas.CheckTest(c, "AWS_VPC_002", checkIfOnlyOneVPC)(s, vpcs, "AWS_VPC_002", &checks)
	yatas.CheckTest(c, "AWS_VPC_003", checkIfOnlyOneGateway)(s, vpcs, "AWS_VPC_003", &checks)
	yatas.CheckTest(c, "AWS_VPC_004", checkIfVPCFLowLogsEnabled)(s, vpcs, "AWS_VPC_004", &checks)
	yatas.CheckTest(c, "AWS_VPC_005", CheckIfAtLeast2Subnets)(s, vpcs, "AWS_VPC_005", &checks)
	yatas.CheckTest(c, "AWS_VPC_006", CheckIfSubnetInDifferentZone)(s, vpcs, "AWS_VPC_006", &checks)
	return checks
}
