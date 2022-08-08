package iam

import (
	"encoding/json"
	"fmt"
	"net/url"

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

func JsonDecodePolicyDocument(policyDocumentJson *string) PolicyDocument {
	// URL Decode the policy document
	var policyDocument PolicyDocument
	decodedValue, err := url.QueryUnescape(*policyDocumentJson)
	if err != nil {
		panic(err)
	}
	// JSON Decode the policy document
	err = json.Unmarshal([]byte(decodedValue), &policyDocument)
	if err != nil {
		panic(err)
	}
	fmt.Println(policyDocument)
	return policyDocument

}
