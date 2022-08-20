package iam

import (
	"context"
	"sync"

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

func GetUserPolicies(users []types.User, s aws.Config) []UserPolicies {
	var wgPolicyForUser sync.WaitGroup
	wgPolicyForUser.Add(len(users))
	queue := make(chan UserPolicies, 10)
	for _, user := range users {
		go GetAllPolicyForUser(&wgPolicyForUser, queue, s, user)
	}
	var userPolicies []UserPolicies
	go func() {
		for user := range queue {
			userPolicies = append(userPolicies, user)
			wgPolicyForUser.Done()
		}

	}()
	wgPolicyForUser.Wait()
	return userPolicies
}

type UserToPoliciesElevate struct {
	UserName string
	Policies [][]string
}

func GetUserToPoliciesElevate(userPolicies []UserPolicies) []UserToPoliciesElevate {
	var usersElevatedPolicies []UserToPoliciesElevate
	for _, user := range userPolicies {
		elevation := CheckPolicyForAllowInRequiredPermission(user.Policies, requiredPermissions)
		if elevation != nil {
			usersElevatedPolicies = append(usersElevatedPolicies, UserToPoliciesElevate{
				UserName: user.UserName,
				Policies: elevation,
			})
		}

	}

	return usersElevatedPolicies
}

func GetAllPolicyForUser(wg *sync.WaitGroup, queueCheck chan UserPolicies, s aws.Config, user types.User) {
	var policyList []Policy
	var wgpolicy sync.WaitGroup
	queue := make(chan *string, 100)
	policies := GetPolicyAttachedToUser(s, user)
	wgpolicy.Add(len(policies))
	for _, policy := range policies {
		go GetPolicyDocument(&wgpolicy, queue, s, policy.PolicyArn)

	}
	go func() {
		for t := range queue {
			policyList = append(policyList, JsonDecodePolicyDocument(t))
			wgpolicy.Done()
		}
	}()
	wgpolicy.Wait()
	queueCheck <- UserPolicies{*user.UserName, policyList}
}

func GetPolicyDocument(wg *sync.WaitGroup, queue chan *string, s aws.Config, policyArn *string) {
	policyVersions := GetAllPolicyVersions(s, policyArn)
	SortPolicyVersions(policyVersions)
	input := &iam.GetPolicyVersionInput{
		PolicyArn: policyArn,
		VersionId: policyVersions[0].VersionId,
	}
	svc := iam.NewFromConfig(s)
	result, err := svc.GetPolicyVersion(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	queue <- result.PolicyVersion.Document
}

func GetPolicyAttachedToUser(s aws.Config, user types.User) []types.AttachedPolicy {
	svc := iam.NewFromConfig(s)
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: user.UserName,
	}
	result, err := svc.ListAttachedUserPolicies(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return result.AttachedPolicies
}

func GetAllPolicyVersions(s aws.Config, policyArn *string) []types.PolicyVersion {
	svc := iam.NewFromConfig(s)
	input := &iam.ListPolicyVersionsInput{
		PolicyArn: policyArn,
	}
	result, err := svc.ListPolicyVersions(context.TODO(), input)
	if err != nil {
		panic(err)
	}

	return result.Versions
}