package iam

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func GetAllUsers(s aws.Config) []types.User {
	svc := iam.NewFromConfig(s)
	input := &iam.ListUsersInput{}
	result, err := svc.ListUsers(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.Users
}

type MFAForUser struct {
	UserName string
	MFAs     []types.MFADevice
}

func GetMfaForUsers(s aws.Config, u []types.User) []MFAForUser {
	svc := iam.NewFromConfig(s)
	input := &iam.ListMFADevicesInput{}
	result, err := svc.ListMFADevices(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var mfaForUsers []MFAForUser
	for _, user := range u {
		mfaForUsers = append(mfaForUsers, MFAForUser{
			UserName: *user.UserName,
			MFAs:     result.MFADevices,
		})
	}
	return mfaForUsers
}

type AccessKeysForUser struct {
	UserName   string
	AccessKeys []types.AccessKeyMetadata
}

func GetAccessKeysForUsers(s aws.Config, u []types.User) []AccessKeysForUser {
	svc := iam.NewFromConfig(s)
	input := &iam.ListAccessKeysInput{}
	result, err := svc.ListAccessKeys(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	var accessKeysForUsers []AccessKeysForUser
	for _, user := range u {
		accessKeysForUsers = append(accessKeysForUsers, AccessKeysForUser{
			UserName:   *user.UserName,
			AccessKeys: result.AccessKeyMetadata,
		})
	}
	return accessKeysForUsers
}
