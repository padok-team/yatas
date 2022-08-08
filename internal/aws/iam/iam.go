package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetAllUsers(s *session.Session) []*iam.User {
	svc := iam.New(s)
	input := &iam.ListUsersInput{}
	result, err := svc.ListUsers(input)
	if err != nil {
		panic(err)
	}
	return result.Users
}

func CheckIf2FAActivated(s *session.Session, users []*iam.User, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "IAM 2FA"
	check.Id = testName
	check.Description = "Check if all users have 2FA activated"
	check.Status = "OK"
	svc := iam.New(s)
	for _, user := range users {
		// List MFA devices for the user
		params := &iam.ListMFADevicesInput{
			UserName: user.UserName,
		}
		resp, err := svc.ListMFADevices(params)
		if err != nil {
			panic(err)
		}
		if len(resp.MFADevices) == 0 {
			check.Status = "FAIL"
			status := "FAIL"
			Message := "2FA is not activated on " + *user.UserName
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		} else {
			status := "OK"
			Message := "2FA is activated on " + *user.UserName
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *user.UserName})
		}
	}
	*c = append(*c, check)
}

func CheckAgeAccessKeyLessThan90Days(s *session.Session, users []*iam.User, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "IAM Access Key Age"
	check.Id = testName
	check.Description = "Check if all users have access key less than 90 days"
	check.Status = "OK"
	svc := iam.New(s)
	for _, user := range users {
		// List access keys for the user
		params := &iam.ListAccessKeysInput{
			UserName: user.UserName,
		}
		resp, err := svc.ListAccessKeys(params)
		if err != nil {
			panic(err)
		}
		now := time.Now()
		for _, accessKey := range resp.AccessKeyMetadata {
			if now.Sub(*accessKey.CreateDate).Hours() > 2160 {
				check.Status = "FAIL"
				status := "FAIL"
				Message := "Access key " + *accessKey.AccessKeyId + " is older than 90 days on " + *user.UserName
				check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *user.UserName})
			} else {
				status := "OK"
				Message := "Access key " + *accessKey.AccessKeyId + " is younger than 90 days on " + *user.UserName
				check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *user.UserName})
			}
		}
	}
	*c = append(*c, check)
}

func RunIAMTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	users := GetAllUsers(s)
	config.CheckTest(c, "AWS_IAM_001", CheckIf2FAActivated)(s, users, "AWS_IAM_001", &checks)
	config.CheckTest(c, "AWS_IAM_002", CheckAgeAccessKeyLessThan90Days)(s, users, "AWS_IAM_002", &checks)
	return checks
}
