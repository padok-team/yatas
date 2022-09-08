package iam

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

func TestSortPolicyVersions(t *testing.T) {
	type args struct {
		policyVersions []types.PolicyVersion
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "SortPolicyVersions",
			args: args{
				policyVersions: []types.PolicyVersion{
					{
						CreateDate: &time.Time{},
					},
					{
						CreateDate: &time.Time{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SortPolicyVersions(tt.args.policyVersions)
			if tt.args.policyVersions[0].CreateDate.After(*tt.args.policyVersions[1].CreateDate) {
				t.Errorf("SortPolicyVersions() = %v", tt.args.policyVersions)
			}
		})
	}
}

func TestJsonDecodePolicyDocument(t *testing.T) {
	type args struct {
		policyDocumentJson *string
	}
	tests := []struct {
		name string
		args args
		want Policy
	}{
		{
			name: "JsonDecodePolicyDocument",
			args: args{
				policyDocumentJson: aws.String("{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}"),
			},
			want: Policy{
				Version: "2012-10-17",
				Statements: []Statement{
					{
						Effect:   "Allow",
						Action:   []string{"*"},
						Resource: []string{"*"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JsonDecodePolicyDocument(tt.args.policyDocumentJson); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonDecodePolicyDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
