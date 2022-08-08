package cloudfront

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func GetAllCloudfront(s *session.Session) []*cloudfront.DistributionSummary {
	svc := cloudfront.New(s)
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(input)
	if err != nil {
		panic(err)
	}
	return result.DistributionList.Items
}

func CheckIfCloudfrontTLS1_2Minimum(s *session.Session, d []*cloudfront.DistributionSummary, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "TLS 1.2 Minimum"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have TLS 1.2 minimum"
	check.Status = "OK"
	for _, cloudfront := range d {
		if cloudfront.ViewerCertificate != nil && strings.Contains(*cloudfront.ViewerCertificate.MinimumProtocolVersion, "TLSv1.2") {
			check.Status = "OK"
			status := "OK"
			Message := "TLS 1.2 minimum is set on " + *cloudfront.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		} else {
			status := "FAIL"
			Message := "TLS 1.2 minimum is not set on " + *cloudfront.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		}
	}
	*c = append(*c, check)
}

func RunCloudFrontTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	d := GetAllCloudfront(s)
	config.CheckTest(c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(s, d, "AWS_CFT_001", &checks)
	return checks
}
