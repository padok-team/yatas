package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/stangirard/yatas/aws/acm"
	"github.com/stangirard/yatas/aws/apigateway"
	"github.com/stangirard/yatas/aws/autoscaling"
	"github.com/stangirard/yatas/aws/cloudfront"
	"github.com/stangirard/yatas/aws/cloudtrail"
	"github.com/stangirard/yatas/aws/dynamodb"
	"github.com/stangirard/yatas/aws/ec2"
	"github.com/stangirard/yatas/aws/ecr"
	"github.com/stangirard/yatas/aws/eks"
	"github.com/stangirard/yatas/aws/guardduty"
	"github.com/stangirard/yatas/aws/iam"
	"github.com/stangirard/yatas/aws/lambda"
	"github.com/stangirard/yatas/aws/loadbalancers"
	"github.com/stangirard/yatas/aws/rds"
	"github.com/stangirard/yatas/aws/s3"
	"github.com/stangirard/yatas/aws/volumes"
	"github.com/stangirard/yatas/aws/vpc"
	"github.com/stangirard/yatas/example"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

// Create a new session that the SDK will use to load
// credentials from. With either SSO or credentials
func initAuth(a yatas.AWS_Account) aws.Config {

	s := initSession(a)
	return s

}

// Create a new session that the SDK will use to load
// credentials from credentials
func createSessionWithCredentials(c yatas.AWS_Account) aws.Config {

	if c.Profile == "" {
		s, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s
	} else {
		s, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Region),
			config.WithSharedConfigProfile(c.Profile),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s
	}

}

// Create a new session that the SDK will use to load
// credentials from the shared credentials file.
// Usefull for SSO
func createSessionWithSSO(c yatas.AWS_Account) aws.Config {

	if c.Profile == "" {
		s, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(c.Region),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s
	} else {
		s, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(c.Region),
			config.WithSharedConfigProfile(c.Profile),
			config.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), 10)
			}),
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s

	}

}

// Create a new session that the SDK will use to load
// credentials from. With either SSO or credentials
func initSession(c yatas.AWS_Account) aws.Config {

	if c.SSO {
		logger.Debug("Using AWS SSO")
		return createSessionWithSSO(c)
	} else {
		logger.Debug("Using AWS credentials")
		return createSessionWithCredentials(c)
	}
}

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

// Here is a real implementation of Greeter
type YatasPlugin struct {
	logger hclog.Logger
}

func (g *YatasPlugin) Run(c *yatas.Config) []yatas.Tests {
	g.logger.Debug("message from YatasPlugin.Run")

	var checksAll []yatas.Tests

	checks, err := runPlugins(c, "aws")
	if err != nil {
		g.logger.Error("Error running plugins", "error", err)
	}
	checksAll = append(checksAll, checks...)
	return checksAll
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	yatasPlugin := &YatasPlugin{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"aws": &example.GreeterPlugin{Impl: yatasPlugin},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

// Run the plugins that are enabled in the config with a switch based on the name of the plugin
func runPlugins(c *yatas.Config, plugin string) ([]yatas.Tests, error) {
	var checksAll []yatas.Tests

	logger.Debug(fmt.Sprint("Running plugin: ", plugin))

	checksAll, err := Run(c)
	if err != nil {
		return nil, err
	}

	return checksAll, nil
}
