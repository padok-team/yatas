package cloudfront

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func GetAllCloudfront(s aws.Config) []types.DistributionSummary {
	svc := cloudfront.NewFromConfig(s)
	input := &cloudfront.ListDistributionsInput{}
	result, err := svc.ListDistributions(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.DistributionList.Items
}

func CheckIfCloudfrontTLS1_2Minimum(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "TLS 1.2 Minimum"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have TLS 1.2 minimum"
	check.Status = "OK"
	for _, cloudfront := range d {
		if cloudfront.ViewerCertificate != nil && strings.Contains(string(cloudfront.ViewerCertificate.MinimumProtocolVersion), "TLSv1.2") {
			check.Status = "OK"
			status := "OK"
			Message := "TLS 1.2 minimum is set on " + *cloudfront.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		} else {
			status := "FAIL"
			Message := "TLS 1.2 minimum is not set on " + *cloudfront.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfHTTPSOnly(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Cloudfront HTTPS Only"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions are HTTPS only"
	check.Status = "OK"
	for _, cloudfront := range d {
		if cloudfront.DefaultCacheBehavior != nil && cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "https-only" || cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "redirect-to-https" {
			check.Status = "OK"
			status := "OK"
			Message := "Cloudfront distribution is HTTPS only on " + *cloudfront.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		} else {
			status := "FAIL"
			Message := "Cloudfront distribution is not HTTPS only on " + *cloudfront.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cloudfront.Id})
		}
	}

	*c = append(*c, check)
	wg.Done()
}

func CheckIfStandardLogginEnabled(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Standard Logging Enabled"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have standard logging enabled"
	check.Status = "OK"
	svc := cloudfront.NewFromConfig(s)
	for _, cc := range d {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		if result.DistributionConfig.Logging != nil && *result.DistributionConfig.Logging.Enabled {
			check.Status = "OK"
			status := "OK"
			Message := "Standard logging is enabled on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		} else {
			status := "FAIL"
			Message := "Standard logging is not enabled on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfCookieLogginEnabled(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "Cookie Logging Enabled"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have cookie logging enabled"
	check.Status = "OK"
	svc := cloudfront.NewFromConfig(s)
	for _, cc := range d {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		if result.DistributionConfig.Logging != nil && *result.DistributionConfig.Logging.Enabled && *result.DistributionConfig.Logging.IncludeCookies {
			check.Status = "OK"
			status := "OK"
			Message := "Cookie logging is enabled on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		} else {
			status := "FAIL"
			Message := "Cookie logging is not enabled on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func CheckIfACLUsed(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, c *[]results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.Name = "ACL Used"
	check.Id = testName
	check.Description = "Check if all cloudfront distributions have an ACL used"
	check.Status = "OK"
	svc := cloudfront.NewFromConfig(s)
	for _, cc := range d {
		input := &cloudfront.GetDistributionConfigInput{
			Id: cc.Id,
		}
		result, err := svc.GetDistributionConfig(context.TODO(), input)
		if err != nil {
			panic(err)
		}
		if *result.DistributionConfig.WebACLId != "" {
			check.Status = "OK"
			status := "OK"
			Message := "ACL is used on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		} else {
			status := "FAIL"
			Message := "ACL is not used on " + *cc.Id
			check.Results = append(check.Results, results.Result{Status: status, Message: Message, ResourceID: *cc.Id})
		}
	}
	*c = append(*c, check)
	wg.Done()
}

func RunChecks(s aws.Config, c *yatas.Config) []results.Check {
	var checks []results.Check
	d := GetAllCloudfront(s)
	var wg sync.WaitGroup

	go yatas.CheckTest(&wg, c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(&wg, s, d, "AWS_CFT_001", &checks)
	go yatas.CheckTest(&wg, c, "AWS_CFT_002", CheckIfHTTPSOnly)(&wg, s, d, "AWS_CFT_002", &checks)
	go yatas.CheckTest(&wg, c, "AWS_CFT_003", CheckIfStandardLogginEnabled)(&wg, s, d, "AWS_CFT_003", &checks)
	go yatas.CheckTest(&wg, c, "AWS_CFT_004", CheckIfCookieLogginEnabled)(&wg, s, d, "AWS_CFT_004", &checks)
	go yatas.CheckTest(&wg, c, "AWS_CFT_005", CheckIfACLUsed)(&wg, s, d, "AWS_CFT_005", &checks)
	wg.Wait()
	return checks
}
