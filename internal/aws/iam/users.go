package iam

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func GetPolicyAttachedToUser(s *session.Session, user *iam.User) []*iam.AttachedPolicy {
	svc := iam.New(s)
	input := &iam.ListAttachedUserPoliciesInput{
		UserName: user.UserName,
	}
	result, err := svc.ListAttachedUserPolicies(input)
	if err != nil {
		panic(err)
	}
	return result.AttachedPolicies
}

func GetAllPolicyVersions(s *session.Session, policyArn *string) []*iam.PolicyVersion {
	svc := iam.New(s)
	input := &iam.ListPolicyVersionsInput{
		PolicyArn: policyArn,
	}
	result, err := svc.ListPolicyVersions(input)
	if err != nil {
		panic(err)
	}

	return result.Versions
}

func SortPolicyVersions(policyVersions []*iam.PolicyVersion) {
	for i := 0; i < len(policyVersions); i++ {
		for j := i + 1; j < len(policyVersions); j++ {
			if policyVersions[i].CreateDate.After(*policyVersions[j].CreateDate) {
				policyVersions[i], policyVersions[j] = policyVersions[j], policyVersions[i]
			}
		}
	}
}

func GetPolicyDocument(s *session.Session, policyArn *string) *string {
	policyVersions := GetAllPolicyVersions(s, policyArn)
	SortPolicyVersions(policyVersions)
	input := &iam.GetPolicyVersionInput{
		PolicyArn: policyArn,
		VersionId: policyVersions[0].VersionId,
	}
	svc := iam.New(s)
	result, err := svc.GetPolicyVersion(input)
	if err != nil {
		panic(err)
	}
	return result.PolicyVersion.Document
}

func JsonDecodePolicyDocument(policyDocumentJson *string) Policy {
	// URL Decode the policy document
	var policyDocument Policy
	decodedValue, _ := url.QueryUnescape(*policyDocumentJson)
	policyDocument.UnmarshalJSON([]byte(decodedValue))
	return policyDocument

}

func GetAllPolicyForUser(s *session.Session, user *iam.User) []Policy {
	var policyList []Policy
	for _, policy := range GetPolicyAttachedToUser(s, user) {
		policyList = append(policyList, JsonDecodePolicyDocument(GetPolicyDocument(s, policy.PolicyArn)))
	}
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
