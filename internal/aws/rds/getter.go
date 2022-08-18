package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
)

func GetListRDS(s aws.Config) []types.DBInstance {
	logger.Debug("Getting list of RDS instances")
	svc := rds.NewFromConfig(s)

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.DBInstances))
	return resp.DBInstances
}
