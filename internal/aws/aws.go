package aws

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stangirard/yatas/internal/aws/apigateway"
	"github.com/stangirard/yatas/internal/aws/autoscaling"
	"github.com/stangirard/yatas/internal/aws/cloudfront"
	"github.com/stangirard/yatas/internal/aws/cloudtrail"
	"github.com/stangirard/yatas/internal/aws/dynamodb"
	"github.com/stangirard/yatas/internal/aws/ec2"
	"github.com/stangirard/yatas/internal/aws/ecr"
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

func Run(c *yatas.Config) ([]results.Check, error) {
	s := initAuth(c)
	logger.Info("Launching AWS checks")
	checks := initTest(s, c)
	return checks, nil
}

func initTest(s aws.Config, c *yatas.Config) []results.Check {

	var checks []results.Check
	var wg sync.WaitGroup

	queue := make(chan []results.Check, 1)
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

	go func() {
		// defer wg.Done() <- Never gets called since the 100 `Done()` calls are made above, resulting in the `Wait()` to continue on before this is executed
		for t := range queue {
			checks = append(checks, t...)
			wg.Done() // ** move the `Done()` call here
		}
	}()
	wg.Wait()

	logger.Info("AWS checks completed ✅")

	return checks
}
