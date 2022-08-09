package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func initAuth(config *yatas.Config) *session.Session {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	s := initSession(config)
	return s

}

func createSessionWithCredentials(c *yatas.Config) *session.Session {
	// Create a new session that the SDK will use to load
	// credentials from credentials
	var s *session.Session
	if c.AWS.Account.Profile == "" {
		s = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(c.AWS.Account.Region),
			}}))
	} else {
		s = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: aws.String(c.AWS.Account.Region),
			},
			Profile: c.AWS.Account.Profile,
		}))
	}

	return s
}

func createSessionWithSSO(c *yatas.Config) *session.Session {
	// Create a new session that the SDK will use to load
	// credentials from the shared credentials file.
	// Usefull for SSO
	var s *session.Session
	if c.AWS.Account.Profile == "" {
		s = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable}))
	} else {
		s = session.Must(session.NewSessionWithOptions(session.Options{
			Profile:           c.AWS.Account.Profile,
			SharedConfigState: session.SharedConfigEnable}))

	}
	return s

}

func initSession(c *yatas.Config) *session.Session {
	// Create a new session that the SDK will use to load
	// credentials from. With either SSO or credentials
	if c.AWS.Account.SSO {
		logger.Debug("Using AWS SSO")
		return createSessionWithSSO(c)
	} else {
		logger.Debug("Using AWS credentials")
		return createSessionWithCredentials(c)
	}
}
