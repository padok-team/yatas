package iam

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetAllUsers(s aws.Config) []types.User {
	svc := iam.NewFromConfig(s)
	input := &iam.ListUsersInput{}
	result, err := svc.ListUsers(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Users
}

func CheckIf2FAActivated(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "IAM 2FA"
	check.Id = testName
	check.Description = "Check if all users have 2FA activated"
	check.Status = "OK"
	svc := iam.NewFromConfig(s)
	for _, user := range users {
		// List MFA devices for the user
		params := &iam.ListMFADevicesInput{
			UserName: user.UserName,
		}
		resp, err := svc.ListMFADevices(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		if len(resp.MFADevices) == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "2FA is not activated on " + *user.UserName
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		} else {
			status := "OK"
			Message := "2FA is activated on " + *user.UserName
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckAgeAccessKeyLessThan90Days(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "IAM Access Key Age"
	check.Id = testName
	check.Description = "Check if all users have access key less than 90 days"
	check.Status = "OK"
	svc := iam.NewFromConfig(s)
	for _, user := range users {
		// List access keys for the user
		params := &iam.ListAccessKeysInput{
			UserName: user.UserName,
		}
		resp, err := svc.ListAccessKeys(context.TODO(), params)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		for _, accessKey := range resp.AccessKeyMetadata {
			if now.Sub(*accessKey.CreateDate).Hours() > 2160 {
				check.Status = "FAIL"
				status := "FAIL"
				Message := "Access key " + *accessKey.AccessKeyId + " is older than 90 days on " + *user.UserName
				check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
			} else {
				status := "OK"
				Message := "Access key " + *accessKey.AccessKeyId + " is younger than 90 days on " + *user.UserName
				check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
			}
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfUserCanElevateRights(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "IAM User Can Elevate Rights"
	check.Id = testName
	check.Description = "Check if  users can elevate rights"
	check.Status = "OK"
	for _, user := range users {
		elevation := CheckPolicyForAllowInRequiredPermission(GetAllPolicyForUser(s, users[0]), requiredPermissions)
		if len(elevation) > 0 {
			check.Status = "FAIL"
			status := "FAIL"
			var Message string
			if len(elevation) > 3 {
				Message = "User " + *user.UserName + " can elevate rights with " + fmt.Sprint(elevation[len(elevation)-3:]) + " only last 3 policies"
			} else {
				Message = "User " + *user.UserName + " can elevate rights with " + fmt.Sprint(elevation)
			}

			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		} else {
			status := "OK"
			Message := "User " + *user.UserName + " cannot elevate rights"
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func RunChecks(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	users := GetAllUsers(s)
	var wg sync.WaitGroup

	go yatas.CheckTest(&wg, c, "AWS_IAM_001", CheckIf2FAActivated)(&wg, s, users, "AWS_IAM_001", &checks)
	go yatas.CheckTest(&wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(&wg, s, users, "AWS_IAM_002", &checks)
	go yatas.CheckTest(&wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(&wg, s, users, "AWS_IAM_003", &checks)
	wg.Wait()
	return checks
}
