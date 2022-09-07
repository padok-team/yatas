package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
	"github.com/stangirard/yatas/plugins/aws/acm"
	"github.com/stangirard/yatas/plugins/aws/apigateway"
	"github.com/stangirard/yatas/plugins/aws/autoscaling"
	"github.com/stangirard/yatas/plugins/aws/cloudfront"
	"github.com/stangirard/yatas/plugins/aws/cloudtrail"
	"github.com/stangirard/yatas/plugins/aws/dynamodb"
	"github.com/stangirard/yatas/plugins/aws/ec2"
	"github.com/stangirard/yatas/plugins/aws/ecr"
	"github.com/stangirard/yatas/plugins/aws/eks"
	"github.com/stangirard/yatas/plugins/aws/guardduty"
	"github.com/stangirard/yatas/plugins/aws/iam"
	"github.com/stangirard/yatas/plugins/aws/lambda"
	"github.com/stangirard/yatas/plugins/aws/loadbalancers"
	"github.com/stangirard/yatas/plugins/aws/rds"
	"github.com/stangirard/yatas/plugins/aws/s3"
	"github.com/stangirard/yatas/plugins/aws/volumes"
	"github.com/stangirard/yatas/plugins/aws/vpc"
)

// Public Functin used to run the AWS tests
func Run(c *yatas.Config) ([]yatas.Tests, error) {
	logger.Info("Launching AWS checks")

	var wg sync.WaitGroup
	var queue = make(chan yatas.Tests, 10)
	var checks []yatas.Tests
	wg.Add(len(c.AWS))
	for _, account := range c.AWS {
		go runTestsForAccount(account, c, queue)
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

// For each account we run the tests. We use a queue to store the results and a waitgroup to wait for all the tests to be done. This allows to run all tests asynchronously.
func runTestsForAccount(account yatas.AWS_Account, c *yatas.Config, queue chan yatas.Tests) {
	s := initAuth(account)
	checks := initTest(s, c, account)
	queue <- checks
}

// Main function that launched all the test for a given account. If a new category is added, it needs to be added here.
func initTest(s aws.Config, c *yatas.Config, a yatas.AWS_Account) yatas.Tests {

	var checks yatas.Tests
	checks.Account = a.Name
	var wg sync.WaitGroup
	queue := make(chan []yatas.Check, 100)
	go yatas.CheckMacroTest(&wg, c, acm.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, s3.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, volumes.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, rds.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, vpc.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, cloudtrail.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, ecr.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, lambda.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, dynamodb.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, ec2.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, cloudfront.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, apigateway.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, autoscaling.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, loadbalancers.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, guardduty.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, iam.RunChecks)(&wg, s, c, queue)
	go yatas.CheckMacroTest(&wg, c, eks.RunChecks)(&wg, s, c, queue)

	go func() {
		for t := range queue {

			checks.Checks = append(checks.Checks, t...)

			wg.Done()

		}
	}()
	wg.Wait()

	logger.Info("AWS checks completed ✅")

	return checks
}
