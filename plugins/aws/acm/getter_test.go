package acm

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

type mockACMApi func()

func (a mockACMApi) ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return &acm.ListCertificatesOutput{
		CertificateSummaryList: []types.CertificateSummary{
			{
				CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
			},
		},
	}, nil

}
func (a mockACMApi) DescribeCertificate(ctx context.Context, params *acm.DescribeCertificateInput, optFns ...func(*acm.Options)) (*acm.DescribeCertificateOutput, error) {
	return &acm.DescribeCertificateOutput{
		Certificate: &types.CertificateDetail{
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
			Status:         types.CertificateStatusExpired,
		},
	}, nil
}

func TestGetCertificates(t *testing.T) {
	type args struct {
		svc ACMGetObjectAPI
	}
	tests := []struct {
		name string
		args args
		want []types.CertificateDetail
	}{
		{
			name: "test",
			args: args{
				svc: mockACMApi(func() {}),
			},
			want: []types.CertificateDetail{
				{
					CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
					Status:         types.CertificateStatusExpired,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCertificates(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCertificates() = %v, want %v", got, tt.want)
			}
		})
	}
}
