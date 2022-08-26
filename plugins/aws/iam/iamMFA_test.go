package iam

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIf2FAActivated(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		mfaForUsers []MFAForUser
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all users have 2FA activated",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				mfaForUsers: []MFAForUser{
					{
						UserName: "test",
						MFAs: []types.MFADevice{
							{
								SerialNumber: aws.String("test"),
							},
						},
					},
				},
				testName: "Check if all users have 2FA activated",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIf2FAActivated(tt.args.checkConfig, tt.args.mfaForUsers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIf2FAActivated() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
