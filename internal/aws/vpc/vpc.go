package vpc

import (
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetListVPC(s *session.Session) []*ec2.Vpc {
	svc := ec2.New(s)
	input := &ec2.DescribeVpcsInput{}
	result, err := svc.DescribeVpcs(input)
	if err != nil {
		panic(err)
	}
	return result.Vpcs
}

func checkCIDR20(s *session.Session, vpcs []*ec2.Vpc, c *[]types.Check) {
	logger.Info("Running AWS_VPC_001")
	var check types.Check
	check.Name = "VPC CIDR"
	check.Id = "AWS_VPC_001"
	check.Description = "Check if VPC CIDR is /20 or bigger"
	check.Status = "OK"
	svc := ec2.New(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeVpcsInput{
			VpcIds: []*string{vpc.VpcId},
		}
		resp, err := svc.DescribeVpcs(params)
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
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "VPC CIDR is /20 or bigger on " + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfVPCFLowLogsEnabled(s *session.Session, vpcs []*ec2.Vpc, c *[]types.Check) {
	logger.Info("Running AWS_VPC_002")
	var check types.Check
	check.Name = "VPC Flow Logs"
	check.Id = "AWS_VPC_002"
	check.Description = "Check if VPC Flow Logs are enabled"
	check.Status = "OK"
	svc := ec2.New(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeFlowLogsInput{
			Filter: []*ec2.Filter{
				{
					Name: aws.String("resource-id"),
					Values: []*string{
						vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeFlowLogs(params)
		if err != nil {
			panic(err)
		}
		if len(resp.FlowLogs) == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC Flow Logs are not enabled on " + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "VPC Flow Logs are enabled on " + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfOnlyOneGateway(s *session.Session, vpcs []*ec2.Vpc, c *[]types.Check) {
	logger.Info("Running AWS_VPC_003")
	var check types.Check
	check.Name = "VPC Gateway"
	check.Id = "AWS_VPC_003"
	check.Description = "Check if VPC has only one gateway"
	check.Status = "OK"
	svc := ec2.New(s)
	for _, vpc := range vpcs {
		params := &ec2.DescribeInternetGatewaysInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("attachment.vpc-id"),
					Values: []*string{
						vpc.VpcId,
					},
				},
			},
		}
		resp, err := svc.DescribeInternetGateways(params)
		if err != nil {
			panic(err)
		}
		if len(resp.InternetGateways) > 1 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC has more than one gateway on " + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "VPC has only one gateway on " + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfOnlyOneVPC(s *session.Session, vpcs []*ec2.Vpc, c *[]types.Check) {
	logger.Info("Running AWS_VPC_004")
	var check types.Check
	check.Name = "VPC Only One"
	check.Id = "AWS_VPC_004"
	check.Description = "Check if VPC has only one VPC"
	check.Status = "OK"
	for _, vpc := range vpcs {
		if len(vpcs) > 1 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "VPC Id:" + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "VPC Id:" + *vpc.VpcId
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}

	*c = append(*c, check)
}

func RunVPCTests(s *session.Session) []types.Check {
	var checks []types.Check
	vpcs := GetListVPC(s)
	checkCIDR20(s, vpcs, &checks)
	checkIfVPCFLowLogsEnabled(s, vpcs, &checks)
	checkIfOnlyOneGateway(s, vpcs, &checks)
	checkIfOnlyOneVPC(s, vpcs, &checks)
	return checks
}
