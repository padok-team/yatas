package rds

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/logger"
)

type RDSGetObjectAPI interface {
	DescribeDBInstances(ctx context.Context, input *rds.DescribeDBInstancesInput, optFns ...func(*rds.Options)) (*rds.DescribeDBInstancesOutput, error)
}

func GetListRDS(svc RDSGetObjectAPI) []types.DBInstance {
	logger.Debug("Getting list of RDS instances")

	params := &rds.DescribeDBInstancesInput{}
	resp, err := svc.DescribeDBInstances(context.TODO(), params)
	if err != nil {
		panic(err)
	}

	logger.Debug(fmt.Sprintf("%v", resp.DBInstances))
	return resp.DBInstances
}
