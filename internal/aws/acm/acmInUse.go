package acm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfACMInUse(checkConfig yatas.CheckConfig, certificates []types.CertificateDetail, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ACM certificates are used", "Check if certificate is in use", testName)
	for _, certificate := range certificates {
		if len(certificate.InUseBy) > 0 {
			Message := "Certificate " + *certificate.CertificateArn + " is in use"
			result := results.Result{Status: "OK", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		} else {
			Message := "Certificate " + *certificate.CertificateArn + " is not in use"
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *certificate.CertificateArn}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
