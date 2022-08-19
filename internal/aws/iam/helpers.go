package iam

import (
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func SortPolicyVersions(policyVersions []types.PolicyVersion) {
	for i := 0; i < len(policyVersions); i++ {
		for j := i + 1; j < len(policyVersions); j++ {
			if policyVersions[i].CreateDate.After(*policyVersions[j].CreateDate) {
				policyVersions[i], policyVersions[j] = policyVersions[j], policyVersions[i]
			}
		}
	}
}

func JsonDecodePolicyDocument(policyDocumentJson *string) Policy {
	// URL Decode the policy document
	var policyDocument Policy
	decodedValue, _ := url.QueryUnescape(*policyDocumentJson)
	policyDocument.UnmarshalJSON([]byte(decodedValue))
	return policyDocument

}
