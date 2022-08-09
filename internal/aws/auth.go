package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func initAuth(config *yatas.Config) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	s := initSession(config)
	return s

}

func createSessionWithCredentials(c *yatas.Config) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from credentials
	if c.AWS.Account.Profile == "" {
		s, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.AWS.Account.Region),
		)
		if err != nil {
			panic(err)
		}
		return s
	} else {
		s, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.AWS.Account.Region),
			config.WithSharedConfigProfile(c.AWS.Account.Profile),
		)
		if err != nil {
			panic(err)
		}
		return s
	}

}

func initSession(c *yatas.Config) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	logger.Debug("Using AWS credentials")
	return createSessionWithCredentials(c)
}
