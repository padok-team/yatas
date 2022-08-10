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
	"github.com/stangirard/yatas/internal/logger"
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

func GetPolicyDocument(s aws.Config, policyArn *string) *string {
	logger.Debug("Getting policy document")
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
	logger.Debug("Got policy document")
	return result.PolicyVersion.Document
}

func JsonDecodePolicyDocument(policyDocumentJson *string) Policy {
	// URL Decode the policy document
	logger.Debug("Decoding policy document")
	var policyDocument Policy
	decodedValue, _ := url.QueryUnescape(*policyDocumentJson)
	policyDocument.UnmarshalJSON([]byte(decodedValue))
	logger.Debug("Decoded policy document")

	return policyDocument

}

func GetAllPolicyForUser(s aws.Config, user types.User) []Policy {
	var policyList []Policy
	queue := make(chan Policy)
	wg := sync.WaitGroup{}
	for _, policy := range GetPolicyAttachedToUser(s, user) {
		go newFunction(&wg, queue, s, policy)
	}
	go func() {
		for policy := range queue {
			policyList = append(policyList, policy)
			wg.Done()
		}
	}()
	wg.Wait()
	return policyList
}

func newFunction(wg *sync.WaitGroup, queue chan Policy, s aws.Config, policy types.AttachedPolicy) {
	wg.Add(1)
	policyTMP := JsonDecodePolicyDocument(GetPolicyDocument(s, policy.PolicyArn))
	logger.Debug("Got policy")

	queue <- policyTMP
}

func CheckPolicyForAllowInRequiredPermission(wg *sync.WaitGroup, queue chan map[string][][]string, user string, policies []Policy, requiredPermission [][]string) {
	// Extract all allow statements from policy
	wg.Add(1)
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
	queue <- map[string][][]string{user: permissionElevationPossible}
}
