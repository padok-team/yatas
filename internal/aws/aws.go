package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stangirard/yatas/internal/aws/s3"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func Run(c *config.Config) ([]types.Check, error) {
	s := initAuth(c)
	logger.Info("Starting AWS tests")
	checks := initTest(s)
	return checks, nil
}

func initTest(s *session.Session) []types.Check {

	fmt.Println("Ran AWS")
	return s3.RunS3Test(s)
}
