package acm

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfCertificateExpiresIn90Days(t *testing.T) {
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
			name: "Check if certificate expires in 90 days",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				certificates: []types.CertificateDetail{
					{
						CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
						DomainName:     aws.String("example.com"),
						Status:         types.CertificateStatusIssued,
						NotAfter:       aws.Time(time.Now().Add(time.Hour * 24 * 91)),
					},
				},
				testName: "Check if certificate expires in 90 days",
			},
			want: "OK",
		},
		{
			name: "Check if certificate expires in 90 days",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				certificates: []types.CertificateDetail{
					{
						CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/12345678-1234-1234-1234-123456789012"),
						DomainName:     aws.String("example.com"),
						Status:         types.CertificateStatusIssued,
						NotAfter:       aws.Time(time.Now().Add(time.Hour * 24 * 89)),
					},
				},
				testName: "Check if certificate expires in 90 days",
			},
			want: "FAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("%s", tt.args.testName)
			CheckIfCertificateExpiresIn90Days(tt.args.checkConfig, tt.args.certificates, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				t.Logf("%v", tt.args.checkConfig.Queue)
				for check := range tt.args.checkConfig.Queue {
					t.Logf("%v", check)
					if check.Status != tt.want {
						t.Errorf("CheckIfCertificateExpiresIn90Days() = %v, want %v", check.Results[0].Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}

			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
