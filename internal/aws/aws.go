package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/aws/acm"
	"github.com/stangirard/yatas/internal/aws/apigateway"
	"github.com/stangirard/yatas/internal/aws/autoscaling"
	"github.com/stangirard/yatas/internal/aws/cloudfront"
	"github.com/stangirard/yatas/internal/aws/cloudtrail"
	"github.com/stangirard/yatas/internal/aws/dynamodb"
	"github.com/stangirard/yatas/internal/aws/ec2"
	"github.com/stangirard/yatas/internal/aws/ecr"
	"github.com/stangirard/yatas/internal/aws/guardduty"
	"github.com/stangirard/yatas/internal/aws/iam"
	"github.com/stangirard/yatas/internal/aws/lambda"
	"github.com/stangirard/yatas/internal/aws/loadbalancers"
	"github.com/stangirard/yatas/internal/aws/rds"
	"github.com/stangirard/yatas/internal/aws/s3"
	"github.com/stangirard/yatas/internal/aws/volumes"
	"github.com/stangirard/yatas/internal/aws/vpc"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Run(c *yatas.Config) ([]results.Tests, error) {
	logger.Info("Launching AWS checks")
	var wg sync.WaitGroup
	var queue = make(chan results.Tests, 10)
	var checks []results.Tests
	wg.Add(len(c.AWS))
	for _, account := range c.AWS {
		go RunTestsForAccount(account, c, queue)
	}
	go func() {
		for t := range queue {
			checks = append(checks, t)
			wg.Done()
		}
	}()
	wg.Wait()

	return checks, nil
}

func RunTestsForAccount(account yatas.AWS_Account, c *yatas.Config, queue chan results.Tests) {
	s := initAuth(account)
	checks := initTest(s, c, account)
	queue <- checks
}

func initTest(s aws.Config, c *yatas.Config, a yatas.AWS_Account) results.Tests {

	var checks results.Tests
	checks.Account = a.Name
	var wg sync.WaitGroup
	queue := make(chan []results.Check)
	go yatas.CheckMacroTest(&wg, c, s3.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, volumes.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, rds.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, vpc.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, cloudtrail.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, ecr.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, lambda.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, dynamodb.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, ec2.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, iam.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, cloudfront.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, apigateway.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, autoscaling.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, loadbalancers.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, guardduty.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, acm.RunChecks)(&wg, s, c, queue)

	go func() {
		for t := range queue {

			checks.Checks = append(checks.Checks, t...)
			wg.Done()
			if c.Progress != nil {
				c.Progress.Add(1)
				c.Progress.RenderBlank()
			}
		}
	}()
	wg.Wait()

	logger.Info("AWS checks completed ✅")

	return checks
}
