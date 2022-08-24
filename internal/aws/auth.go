package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func initAuth(a yatas.AWS_Account) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	s := initSession(a)
	return s

}

func createSessionWithCredentials(c yatas.AWS_Account) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from credentials
	if c.Profile == "" {
		s, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(c.Region),
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
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s
	}

}

func createSessionWithSSO(c yatas.AWS_Account) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from the shared credentials file.
	// Usefull for SSO
	if c.Profile == "" {
		s, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion(c.Region),
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
			config.WithRetryMode(aws.RetryMode(aws.RetryModeAdaptive)),
		)
		if err != nil {
			panic(err)
		}
		return s

	}

}

func initSession(c yatas.AWS_Account) aws.Config {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	if c.SSO {
		logger.Debug("Using AWS SSO")
		return createSessionWithSSO(c)
	} else {
		logger.Debug("Using AWS credentials")
		return createSessionWithCredentials(c)
	}
}
