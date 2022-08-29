package iam

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfUserLastPasswordUse120Days(checkConfig yatas.CheckConfig, users []types.User, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("IAM Users have not used their password for 120 days", "Check if all users have not used their password for 120 days", testName)
	for _, user := range users {
		if user.PasswordLastUsed != nil {
			if time.Since(*user.PasswordLastUsed).Hours() > 120*24 {
				Message := "Password has not been used for more than 120 days on " + *user.UserName
				result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)
			} else {
				Message := "Password has been used in the last 120 days on " + *user.UserName
				result := yatas.Result{Status: "OK", Message: Message, ResourceID: *user.UserName}
				check.AddResult(result)
			}
		} else {
			Message := "Password has never been used on " + *user.UserName
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *user.UserName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
