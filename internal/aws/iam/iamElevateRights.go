package iam

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfUserCanElevateRights(checkConfig yatas.CheckConfig, userToPolociesElevated []UserToPoliciesElevate, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM User Can Elevate Rights", "Check if  users can elevate rights", testName)
	for _, userPol := range userToPolociesElevated {
		if len(userPol.Policies) > 0 {
			var Message string
			if len(userPol.Policies) > 3 {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(userPol.Policies[len(userPol.Policies)-3:]) + " only last 3 policies"
			} else {
				Message = "User " + userPol.UserName + " can elevate rights with " + fmt.Sprint(userPol.Policies)
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
