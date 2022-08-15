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

func CheckIf2FAActivated(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM 2FA", "Check if all users have 2FA activated", testName)
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
			Message := "2FA is not activated on " + *user.UserName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
			check.AddResult(result)
		} else {
			Message := "2FA is activated on " + *user.UserName
			result := results.Result{Status: "OK", Message: Message, ResourceID: *user.UserName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckAgeAccessKeyLessThan90Days(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM Access Key Age", "Check if all users have access key less than 90 days", testName)
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
				Message := "Access key " + *accessKey.AccessKeyId + " is older than 90 days on " + *user.UserName
				result := results.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)

			} else {
				Message := "Access key " + *accessKey.AccessKeyId + " is younger than 90 days on " + *user.UserName
				result := results.Result{Status: "OK", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)
			}
		}
	}
	queueToAdd <- check
}

type UserPolicies struct {
	UserName string
	Policies []Policy
}

func CheckIfUserCanElevateRights(wg *sync.WaitGroup, s aws.Config, users []types.User, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM User Can Elevate Rights", "Check if  users can elevate rights", testName)
	var wgPolicyForUser sync.WaitGroup
	queue := make(chan UserPolicies, len(users))
	wgPolicyForUser.Add(len(users))
	for _, user := range users {
		go GetAllPolicyForUser(&wgPolicyForUser, queue, s, user)
	}
	var userPolicies []UserPolicies
	go func() {
		for user := range queue {
			userPolicies = append(userPolicies, user)
			wgPolicyForUser.Done()
		}

	}()
	wgPolicyForUser.Wait()
	for _, userPol := range userPolicies {
		elevation := CheckPolicyForAllowInRequiredPermission(userPol.Policies, requiredPermissions)
		if len(elevation) > 0 {
			var Message string
			if len(elevation) > 3 {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(elevation[len(elevation)-3:]) + " only last 3 policies"
			} else {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(elevation)
			}
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: userPol.UserName}
			check.AddResult(result)

		} else {
			Message := "User " + userPol.UserName + " cannot elevate rights"
			result := results.Result{Status: "OK", Message: Message, ResourceID: userPol.UserName}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checks []results.Check
	users := GetAllUsers(s)
	var wg sync.WaitGroup
	queueResults := make(chan results.Check, 10)

	go yatas.CheckTest(&wg, c, "AWS_IAM_001", CheckIf2FAActivated)(&wg, s, users, "AWS_IAM_001", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(&wg, s, users, "AWS_IAM_002", queueResults)
	go yatas.CheckTest(&wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(&wg, s, users, "AWS_IAM_003", queueResults)
	go func() {
		for t := range queueResults {
			checks = append(checks, t)
			wg.Done()
		}
	}()

	wg.Wait()

	queue <- checks
}
