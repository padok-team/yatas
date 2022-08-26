package acm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfACMValid(checkConfig yatas.CheckConfig, certificates []types.CertificateDetail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ACM certificates are valid", "Check if certificate is valid", testName)
	for _, certificate := range certificates {
		if certificate.Status == types.CertificateStatusIssued || certificate.Status == types.CertificateStatusInactive {
			Message := "Certificate " + *certificate.CertificateArn + " is valid"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not valid"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
