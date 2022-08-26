package iam

import (
	"fmt"
	"time"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckAgeAccessKeyLessThan90Days(checkConfig yatas.CheckConfig, accessKeysForUsers []AccessKeysForUser, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("IAM access key younger than 90 days", "Check if all users have access key less than 90 days", testName)
	for _, accesskeyforuser := range accessKeysForUsers {
		now := time.Now()
		for _, accessKey := range accesskeyforuser.AccessKeys {
			if now.Sub(*accessKey.CreateDate).Hours() > 2160 {
				Message := "Access key " + *accessKey.AccessKeyId + " is older than 90 days on " + accesskeyforuser.UserName
				result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: accesskeyforuser.UserName}
				check.AddResult(result)

			} else {
				Message := "Access key " + *accessKey.AccessKeyId + " is younger than 90 days on " + accesskeyforuser.UserName
				result := yatas.Result{Status: "OK", Message: Message, ResourceID: accesskeyforuser.UserName}
				check.AddResult(result)
			}
		}
	}
	checkConfig.Queue <- check
}
