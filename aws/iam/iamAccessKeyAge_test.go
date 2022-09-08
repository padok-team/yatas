package iam

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckAgeAccessKeyLessThan90Days(t *testing.T) {
	type args struct {
		checkConfig        yatas.CheckConfig
		accessKeysForUsers []AccessKeysForUser
		testName           string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all users have access key less than 90 days",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				accessKeysForUsers: []AccessKeysForUser{
					{
						UserName: "test",
						AccessKeys: []types.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("test"),
								CreateDate:  aws.Time(time.Now()),
							},
						},
					},
				},
				testName: "Check if all users have access key less than 90 days",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckAgeAccessKeyLessThan90Days(tt.args.checkConfig, tt.args.accessKeysForUsers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckAgeAccessKeyLessThan90Days() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckAgeAccessKeyLessThan90DaysFail(t *testing.T) {
	type args struct {
		checkConfig        yatas.CheckConfig
		accessKeysForUsers []AccessKeysForUser
		testName           string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all users have access key less than 90 days",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				accessKeysForUsers: []AccessKeysForUser{
					{
						UserName: "test",
						AccessKeys: []types.AccessKeyMetadata{
							{
								AccessKeyId: aws.String("test"),
								CreateDate:  aws.Time(time.Now().Add(-time.Hour * 24 * 91)),
							},
						},
					},
				},
				testName: "Check if all users have access key less than 90 days",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckAgeAccessKeyLessThan90Days(tt.args.checkConfig, tt.args.accessKeysForUsers, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckAgeAccessKeyLessThan90Days() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
