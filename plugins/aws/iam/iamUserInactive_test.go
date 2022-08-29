package iam

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfUserLastPasswordUse120Days(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		users       []types.User
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestCheckIfUserLastPasswordUse120Days",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				users:       []types.User{},
				testName:    "AWS_IAM_001",
			},
			want: "OK",
		},
		{
			name: "TestCheckIfUserLastPasswordUse120Days",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				users: []types.User{
					{
						PasswordLastUsed: nil,
						UserName:         aws.String("test"),
					},
				},
				testName: "AWS_IAM_001",
			},
			want: "FAIL",
		},
		{
			name: "TestCheckIfUserLastPasswordUse120Days",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				users: []types.User{
					{
						PasswordLastUsed: aws.Time(time.Now().Add(-121 * 24 * time.Hour)),
						UserName:         aws.String("test"),
					},
				},
				testName: "AWS_IAM_001",
			},
			want: "FAIL",
		},
		{
			name: "TestCheckIfUserLastPasswordUse120Days",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				users: []types.User{
					{
						PasswordLastUsed: aws.Time(time.Now().Add(-20 * 24 * time.Hour)),
						UserName:         aws.String("test"),
					},
				},
				testName: "AWS_IAM_001",
			},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfUserLastPasswordUse120Days(tt.args.checkConfig, tt.args.users, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfUserLastPasswordUse120Days() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
