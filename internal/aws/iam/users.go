package iam

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

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

func SortPolicyVersions(policyVersions []types.PolicyVersion) {
	for i := 0; i < len(policyVersions); i++ {
		for j := i + 1; j < len(policyVersions); j++ {
			if policyVersions[i].CreateDate.After(*policyVersions[j].CreateDate) {
				policyVersions[i], policyVersions[j] = policyVersions[j], policyVersions[i]
			}
		}
	}
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

func JsonDecodePolicyDocument(policyDocumentJson *string) Policy {
	// URL Decode the policy document
	var policyDocument Policy
	decodedValue, _ := url.QueryUnescape(*policyDocumentJson)
	policyDocument.UnmarshalJSON([]byte(decodedValue))
	return policyDocument

}

func GetAllPolicyForUser(s aws.Config, user types.User) []Policy {
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
	return policyList
}

func CheckPolicyForAllowInRequiredPermission(policies []Policy, requiredPermission [][]string) [][]string {
	// Extract all allow statements from policy
	allowStatements := make([]Statement, 0)
	for _, policy := range policies {
		for _, statement := range policy.Statements {
			if statement.Effect == "Allow" {
				allowStatements = append(allowStatements, statement)
			}
		}
	}
	var permissionElevationPossible = [][]string{}
	// Check if any statement is in requiredPermissions
	for _, permissions := range requiredPermissions {
		// Create a map of permissions and false
		permissionMap := make(map[string]bool)
		for _, permission := range permissions {
			permissionMap[permission] = false
		}
		for _, permission := range permissions {
			for _, statement := range allowStatements {
				for _, actions := range statement.Action {
					actions = strings.ReplaceAll(actions, "*", ".*")
					// If regex actions matches permission actions, return true
					found, err := regexp.MatchString(actions, permission)
					if err != nil {
						panic(err)
					}
					if found {
						permissionMap[permission] = true
					}
				}
			}
		}
		// If all permissions are true, return true
		permissionsBool := true
		for _, permission := range permissionMap {
			if !permission {
				permissionsBool = false
			}
		}
		if permissionsBool {
			permissionElevationPossible = append(permissionElevationPossible, permissions)
		}
	}

	return permissionElevationPossible
}
