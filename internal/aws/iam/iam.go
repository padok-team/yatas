package iam

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfUserCanElevateRights(checkConfig yatas.CheckConfig, users []types.User, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM User Can Elevate Rights", "Check if  users can elevate rights", testName)
	var wgPolicyForUser sync.WaitGroup
	queue := make(chan UserPolicies, len(users))
	wgPolicyForUser.Add(len(users))
	for _, user := range users {
		go GetAllPolicyForUser(&wgPolicyForUser, queue, checkConfig.ConfigAWS, user)
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
	checkConfig.Queue <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	users := GetAllUsers(s)
	mfaForUsers := GetMfaForUsers(s, users)
	accessKeysForUsers := GetAccessKeysForUsers(s, users)

	go yatas.CheckTest(checkConfig.Wg, c, "AWS_IAM_001", CheckIf2FAActivated)(checkConfig, mfaForUsers, "AWS_IAM_001")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(checkConfig, accessKeysForUsers, "AWS_IAM_002")
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_IAM_003", CheckIfUserCanElevateRights)(checkConfig, users, "AWS_IAM_003")
	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
