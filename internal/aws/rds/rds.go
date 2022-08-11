package rds

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetListRDS(s aws.Config) []types.DBInstance {
	logger.Debug("Getting list of RDS instances")
	svc := rds.NewFromConfig(s)

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.DBInstances))
	return resp.DBInstances
}

func checkIfEncryptionEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Encryption"
	check.Id = testName
	check.Description = "Check if RDS encryption is enabled"
	check.Status = "OK"
	svc := rds.NewFromConfig(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if !resp.DBInstances[0].StorageEncrypted {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS encryption is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		} else {
			status := "OK"
			Message := "RDS encryption is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func checkIfBackupEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Backup"
	check.Id = testName
	check.Description = "Check if RDS backup is enabled"
	check.Status = "OK"
	svc := rds.NewFromConfig(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.DBInstances[0].BackupRetentionPeriod == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS backup is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		} else {
			status := "OK"
			Message := "RDS backup is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func checkIfAutoUpgradeEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Minor Auto Upgrade"
	check.Id = testName
	check.Description = "Check if RDS minor auto upgrade is enabled"
	check.Status = "OK"
	svc := rds.NewFromConfig(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if !resp.DBInstances[0].AutoMinorVersionUpgrade {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS auto upgrade is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		} else {
			status := "OK"
			Message := "RDS auto upgrade is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func checkIfRDSPrivateEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Private"
	check.Id = testName
	check.Description = "Check if RDS private is enabled"
	check.Status = "OK"
	svc := rds.NewFromConfig(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.DBInstances[0].PubliclyAccessible {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS private is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		} else {
			status := "OK"
			Message := "RDS private is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func CheckIfLoggingEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Logging"
	check.Id = testName
	check.Description = "Check if RDS logging is enabled"
	check.Status = "OK"
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					status := "OK"
					Message := "RDS logging is enabled on " + *instance.DBInstanceIdentifier
					check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
				}
			}
		} else {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func CheckIfDeleteProtectionEnabled(wg *sync.WaitGroup, s aws.Config, instances []types.DBInstance, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "RDS Delete Protection"
	check.Id = testName
	check.Description = "Check if RDS delete protection is enabled"
	check.Status = "OK"
	for _, instance := range instances {
		if instance.DeletionProtection {
			status := "OK"
			Message := "RDS delete protection is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		} else {
			status := "FAIL"
			Message := "RDS delete protection is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *instance.DBInstanceArn})
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	instances := GetListRDS(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)

	go yatas.CheckTest(&wg, c, "AWS_RDS_001", checkIfEncryptionEnabled)(&wg, s, instances, "AWS_RDS_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_RDS_002", checkIfBackupEnabled)(&wg, s, instances, "AWS_RDS_002", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(&wg, s, instances, "AWS_RDS_003", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(&wg, s, instances, "AWS_RDS_004", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_RDS_005", CheckIfLoggingEnabled)(&wg, s, instances, "AWS_RDS_005", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_RDS_006", CheckIfDeleteProtectionEnabled)(&wg, s, instances, "AWS_RDS_006", queueResults)
	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
