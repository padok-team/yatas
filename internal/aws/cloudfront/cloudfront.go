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

func CheckIfCloudfrontTLS1_2Minimum(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("TLS 1.2 Minimum", "Check if all cloudfront distributions have TLS 1.2 minimum", testName)
	for _, cloudfront := range d {
		if cloudfront.ViewerCertificate != nil && strings.Contains(string(cloudfront.ViewerCertificate.MinimumProtocolVersion), "TLSv1.2") {
			Message := "TLS 1.2 minimum is set on " + *cloudfront.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		} else {
			Message := "TLS 1.2 minimum is not set on " + *cloudfront.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfHTTPSOnly(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cloudfront HTTPS Only", "Check if all cloudfront distributions are HTTPS only", testName)
	for _, cloudfront := range d {
		if cloudfront.DefaultCacheBehavior != nil && cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "https-only" || cloudfront.DefaultCacheBehavior.ViewerProtocolPolicy == "redirect-to-https" {
			Message := "Cloudfront distribution is HTTPS only on " + *cloudfront.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		} else {
			Message := "Cloudfront distribution is not HTTPS only on " + *cloudfront.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cloudfront.Id}
			check.AddResult(result)
		}
	}

	queueToAdd <- check
}

func CheckIfStandardLogginEnabled(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Standard Logging Enabled", "Check if all cloudfront distributions have standard logging enabled", testName)
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
			Message := "Standard logging is enabled on " + *cc.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		} else {
			Message := "Standard logging is not enabled on " + *cc.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfCookieLogginEnabled(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("Cookies Logging Enabled", "Check if all cloudfront distributions have cookies logging enabled", testName)
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
			Message := "Cookie logging is enabled on " + *cc.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		} else {
			Message := "Cookie logging is not enabled on " + *cc.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func CheckIfACLUsed(wg *sync.WaitGroup, s aws.Config, d []types.DistributionSummary, testName string, queueToAdd chan results.Check) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check results.Check
	check.InitCheck("ACL Used", "Check if all cloudfront distributions have an ACL used", testName)
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
			Message := "ACL is used on " + *cc.Id
			result := results.Result{Status: "OK", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		} else {
			Message := "ACL is not used on " + *cc.Id
			result := results.Result{Status: "FAIL", Message: Message, ResourceID: *cc.Id}
			check.AddResult(result)
		}
	}
	queueToAdd <- check
}

func RunChecks(wa *sync.WaitGroup, s aws.Config, c *yatas.Config, queue chan []results.Check) {

	var checkConfig yatas.CheckConfig
	checkConfig.Init(s, c)
	var checks []results.Check
	d := GetAllCloudfront(s)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_001", CheckIfCloudfrontTLS1_2Minimum)(checkConfig.Wg, checkConfig.ConfigAWS, d, "AWS_CFT_001", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_002", CheckIfHTTPSOnly)(checkConfig.Wg, checkConfig.ConfigAWS, d, "AWS_CFT_002", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_003", CheckIfStandardLogginEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, d, "AWS_CFT_003", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_004", CheckIfCookieLogginEnabled)(checkConfig.Wg, checkConfig.ConfigAWS, d, "AWS_CFT_004", checkConfig.Queue)
	go yatas.CheckTest(checkConfig.Wg, c, "AWS_CFT_005", CheckIfACLUsed)(checkConfig.Wg, checkConfig.ConfigAWS, d, "AWS_CFT_005", checkConfig.Queue)

	go func() {
		for t := range checkConfig.Queue {
			checks = append(checks, t)
			checkConfig.Wg.Done()
		}
	}()

	checkConfig.Wg.Wait()

	queue <- checks
}
