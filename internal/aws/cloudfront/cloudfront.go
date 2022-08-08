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

func CheckIfHTTPSOnly(s *session.Session, d []*cloudfront.DistributionSummary, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Cloudfront HTTPS Only"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions are HTTPS only"
	check.Status = "OK"
	for _, cloudfront := range d {
		if cloudfront.DefaultCacheBehavior != nil && *cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "https-only" || *cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "redirect-to-https" {
			check.Status = "OK"
			status := "OK"
			Message := "Cloudfront distribution is HTTPS only on " + *cloudfront.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		} else {
			status := "FAIL"
			Message := "Cloudfront distribution is not HTTPS only on " + *cloudfront.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		}
	}

	*c = append(*c, check)
}

func CheckIfStandardLogginEnabled(s *session.Session, d []*cloudfront.DistributionSummary, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Standard Logging Enabled"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have standard logging enabled"
	check.Status = "OK"
	svc := cloudfront.New(s)
	for _, cc := range d {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(input)
		if err != nil {
			panic(err)
		}
		if result.DistributionConfig.Logging != nil && *result.DistributionConfig.Logging.Enabled {
			check.Status = "OK"
			status := "OK"
			Message := "Standard logging is enabled on " + *cc.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		} else {
			status := "FAIL"
			Message := "Standard logging is not enabled on " + *cc.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		}
	}
	*c = append(*c, check)
}

func CheckIfCookieLogginEnabled(s *session.Session, d []*cloudfront.DistributionSummary, testName string, c *[]types.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check types.Check
	check.Name = "Cookie Logging Enabled"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have cookie logging enabled"
	check.Status = "OK"
	svc := cloudfront.New(s)
	for _, cc := range d {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(input)
		if err != nil {
			panic(err)
		}
		if result.DistributionConfig.Logging != nil && *result.DistributionConfig.Logging.Enabled && *result.DistributionConfig.Logging.IncludeCookies {
			check.Status = "OK"
			status := "OK"
			Message := "Cookie logging is enabled on " + *cc.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		} else {
			status := "FAIL"
			Message := "Cookie logging is not enabled on " + *cc.Id
			check.Results = append(check.Results, types.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		}
	}
	*c = append(*c, check)
}

func RunCloudFrontTests(s *session.Session, c *config.Config) []types.Check {
	var checks []types.Check
	d := GetAllCloudfront(s)
	config.CheckTest(c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(s, d, "AWS_CFT_001", &checks)
	config.CheckTest(c, "AWS_CFT_002", CheckIfHTTPSOnly)(s, d, "AWS_CFT_002", &checks)
	config.CheckTest(c, "AWS_CFT_003", CheckIfStandardLogginEnabled)(s, d, "AWS_CFT_003", &checks)
	config.CheckTest(c, "AWS_CFT_004", CheckIfCookieLogginEnabled)(s, d, "AWS_CFT_004", &checks)
	return checks
}
