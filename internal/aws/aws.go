package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stangirard/yatas/internal/aws/rds"
	"github.com/stangirard/yatas/internal/aws/s3"
	"github.com/stangirard/yatas/internal/aws/volumes"
	"github.com/stangirard/yatas/internal/config"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/types"
)

func Run(c *config.Config) ([]types.Check, error) {
	s := initAuth(c)
	logger.Info("Launching AWS checks")
	checks := initTest(s)
	return checks, nil
}

func initTest(s *session.Session) []types.Check {

	var checks []types.Check
	checks = append(checks, s3.RunS3Test(s)...)
	checks = append(checks, volumes.RunVolumesTest(s)...)
	checks = append(checks, rds.RunRDSTests(s)...)
	logger.Info("AWS checks completed ✅")

	return checks
}
