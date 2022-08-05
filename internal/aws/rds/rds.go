package rds

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetListRDS(s *session.Session) []*rds.DBInstance {
	logger.Debug("Getting list of RDS instances")
	svc := rds.New(s)

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.DBInstances))
	return resp.DBInstances
}

func checkIfEncryptionEnabled(s *session.Session, instances []*rds.DBInstance, c *[]types.Check) {
	var check types.Check
	check.Name = "RDS Encryption"
	check.Id = "AWS_RDS_001"
	check.Description = "Check if RDS encryption is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].StorageEncrypted == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS encryption is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS encryption is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfBackupEnabled(s *session.Session, instances []*rds.DBInstance, c *[]types.Check) {
	var check types.Check
	check.Name = "RDS Backup"
	check.Id = "AWS_RDS_002"
	check.Description = "Check if RDS backup is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].BackupRetentionPeriod == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS backup is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS backup is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfAutoUpgradeEnabled(s *session.Session, instances []*rds.DBInstance, c *[]types.Check) {
	var check types.Check
	check.Name = "RDS Minor Auto Upgrade"
	check.Id = "AWS_RDS_003"
	check.Description = "Check if RDS minor auto upgrade is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].AutoMinorVersionUpgrade == false {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS auto upgrade is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS auto upgrade is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func checkIfRDSPrivateEnabled(s *session.Session, instances []*rds.DBInstance, c *[]types.Check) {
	var check types.Check
	check.Name = "RDS Private"
	check.Id = "AWS_RDS_004"
	check.Description = "Check if RDS private is enabled"
	check.Status = "OK"
	svc := rds.New(s)
	for _, instance := range instances {
		params := &rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		}
		resp, err := svc.DescribeDBInstances(params)
		if err != nil {
			panic(err)
		}
		if *resp.DBInstances[0].PubliclyAccessible == true {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "RDS private is not enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		} else {
			status := "OK"
			Message := "RDS private is enabled on " + *instance.DBInstanceIdentifier
			check.Results = append(check.Results, types.Result{Status: status, Message: Message})
		}
	}
	*c = append(*c, check)
}

func RunRDSTests(s *session.Session) []types.Check {
	var checks []types.Check
	instances := GetListRDS(s)
	checkIfEncryptionEnabled(s, instances, &checks)
	checkIfBackupEnabled(s, instances, &checks)
	checkIfAutoUpgradeEnabled(s, instances, &checks)
	checkIfRDSPrivateEnabled(s, instances, &checks)
	return checks
}
