package acm

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfACMValid(t *testing.T) {
	type args struct {
		checkConfig  yatas.CheckConfig
		certificates []types.CertificateDetail
		testName     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check if all ACM certificates are valid",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				certificates: []types.CertificateDetail{
					{
						CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
						Status:         types.CertificateStatusIssued,
					},
				},
				testName: "Check if all ACM certificates are valid",
			},
			want: "OK",
		},
		{
			name: "Check if all ACM certificates are not valid",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				certificates: []types.CertificateDetail{
					{
						CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
						Status:         types.CertificateStatusPendingValidation,
					},
				},
				testName: "Check if all ACM certificates are not valid",
			},
			want: "FAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfACMValid(tt.args.checkConfig, tt.args.certificates, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfACMValid() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
