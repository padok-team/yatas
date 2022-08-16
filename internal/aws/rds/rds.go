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

func checkIfEncryptionEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Encryption", "Check if RDS encryption is enabled", testName)
	svc := rds.NewFromConfig(checkConfig.ConfigAWS)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if !resp.DBInstances[0].StorageEncrypted {
			Message := "RDS encryption is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS encryption is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func checkIfBackupEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Backup", "Check if RDS backup is enabled", testName)
	svc := rds.NewFromConfig(checkConfig.ConfigAWS)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.DBInstances[0].BackupRetentionPeriod == 0 {
			Message := "RDS backup is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS backup is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func checkIfAutoUpgradeEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Minor Auto Upgrade", "Check if RDS minor auto upgrade is enabled", testName)
	svc := rds.NewFromConfig(checkConfig.ConfigAWS)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if !resp.DBInstances[0].AutoMinorVersionUpgrade {
			Message := "RDS auto upgrade is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.Results = append(check.Results, result)
		} else {
			Message := "RDS auto upgrade is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func checkIfRDSPrivateEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Private", "Check if RDS private is enabled", testName)
	svc := rds.NewFromConfig(checkConfig.ConfigAWS)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if resp.DBInstances[0].PubliclyAccessible {
			Message := "RDS private is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS private is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfLoggingEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Logging", "Check if RDS logging is enabled", testName)
	for _, instance := range instances {
		if instance.EnabledCloudwatchLogsExports != nil {
			for _, export := range instance.EnabledCloudwatchLogsExports {
				if export == "audit" {
					Message := "RDS logging is enabled on " + *instance.DBInstanceIdentifier
					result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
					check.AddResult(result)
					break
				}
			}
		} else {
			Message := "RDS logging is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.Results = append(check.Results, result)
		}
	}
	checkConfig.Queue <- check
}

func CheckIfDeleteProtectionEnabled(checkConfig yatas.CheckConfig, instances []types.DBInstance, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("RDS Delete Protection", "Check if RDS delete protection is enabled", testName)
	for _, instance := range instances {
		if instance.DeletionProtection {
			Message := "RDS delete protection is enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "OK", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		} else {
			Message := "RDS delete protection is not enabled on " + *instance.DBInstanceIdentifier
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *instance.DBInstanceArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	instances := GetListRDS(s)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_001", checkIfEncryptionEnabled)(checkConfig, instances, "AWS_RDS_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_002", checkIfBackupEnabled)(checkConfig, instances, "AWS_RDS_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_003", checkIfAutoUpgradeEnabled)(checkConfig, instances, "AWS_RDS_003")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_004", checkIfRDSPrivateEnabled)(checkConfig, instances, "AWS_RDS_004")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_005", CheckIfLoggingEnabled)(checkConfig, instances, "AWS_RDS_005")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_RDS_006", CheckIfDeleteProtectionEnabled)(checkConfig, instances, "AWS_RDS_006")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
