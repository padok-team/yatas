package ec2

import (
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func GetEC2IamInstanceProfile(e types.Instance) string {
	return *e.IamInstanceProfile.Arn
}
