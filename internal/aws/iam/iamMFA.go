package iam

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIf2FAActivated(checkConfig yatas.CheckConfig, mfaForUsers []MFAForUser, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("IAM 2FA", "Check if all users have 2FA activated", testName)
	for _, mfaForUser := range mfaForUsers {
		if len(mfaForUser.MFAs) == 0 {
			Message := "2FA is not activated on " + mfaForUser.UserName
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: mfaForUser.UserName}
			check.AddResult(result)
		} else {
			Message := "2FA is activated on " + mfaForUser.UserName
			result := results.Result{Status: "OK", Message: Message, ResourceID: mfaForUser.UserName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
