package vpc

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

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

func checkCIDR20(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("VPC CIDR", "Check if VPC CIDR is /20 or bigger", testName)
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
			Message := "VPC CIDR is not /20 or bigger on " + *vpc.VpcId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC CIDR is /20 or bigger on " + *vpc.VpcId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func checkIfVPCFLowLogsEnabled(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("VPC Flow Logs", "Check if VPC Flow Logs are enabled", testName)
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
			Message := "VPC Flow Logs are not enabled on " + *vpc.VpcId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC Flow Logs are enabled on " + *vpc.VpcId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func checkIfOnlyOneGateway(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("VPC Gateway", "Check if VPC has only one gateway", testName)
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
			Message := "VPC has more than one gateway on " + *vpc.VpcId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.Results = append(check.Results, result)
		} else {
			Message := "VPC has only one gateway on " + *vpc.VpcId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func checkIfOnlyOneVPC(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("VPC Only One", "Check if VPC has only one VPC", testName)
	for _, vpc := range vpcs {
		if len(vpcs) > 1 {
			Message := "VPC Id:" + *vpc.VpcId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC Id:" + *vpc.VpcId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}

	queueToAdd <- check
}

func CheckIfSubnetInDifferentZone(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Subnets in different zone", "Check if Subnet are in different zone", testName)
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
			Message := "Subnets are in different zone on " + *vpc.VpcId
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.Results = append(check.Results, result)
		} else {
			Message := "Subnets are in same zone on " + *vpc.VpcId
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.Results = append(check.Results, result)
		}
	}
	queueToAdd <- check
}

func CheckIfAtLeast2Subnets(wg *sync.WaitGroup, s aws.Config, vpcs []types.Vpc, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("At least 2 subnets", "Check if VPC has at least 2 subnets", testName)
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
			Message := "VPC " + *vpc.VpcId + " has less than 2 subnets"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		} else {
			Message := "VPC " + *vpc.VpcId + " has at least 2 subnets"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *vpc.VpcId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	vpcs := GetListVPC(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_001", checkCIDR20)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_001", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_002", checkIfOnlyOneVPC)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_002", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_003", checkIfOnlyOneGateway)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_003", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_004", checkIfVPCFLowLogsEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_004", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_005", CheckIfAtLeast2Subnets)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_005", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_VPC_006", CheckIfSubnetInDifferentZone)(checkConfig.Wg, checkConfig.ConfigAWS, vpcs, "AWS_VPC_006", checkConfig.Queue)
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
