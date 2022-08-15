package ec2

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetEC2s(s aws.Config) []types.Instance {
	svc := ec2.NewFromConfig(s)
	input := &ec2.DescribeInstancesInput{}
	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var instances []types.Instance
	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}

	return instances
}

func CheckIfEC2PublicIP(wg *sync.WaitGroup, s aws.Config, instances []types.Instance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("EC2 Public IP", "Check if all instances have a public IP", testName)
	for _, instance := range instances {
		if instance.PublicIpAddress != nil {
			Message := "EC2 instance " + *instance.InstanceId + " has a public IP" + *instance.PublicIpAddress
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		} else {
			Message := "EC2 instance " + *instance.InstanceId + " has no public IP "
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.InstanceId}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	instances := GetEC2s(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)
	go yatas.CheckTest(&wg, c, "AWS_EC2_001", CheckIfEC2PublicIP)(&wg, s, instances, "AWS_EC2_001", queueResults)

	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
