package iam

import (
	"sync"
	"testing"

	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfUserCanElevateRights(t *testing.T) {
	type args struct {
		checkConfig            yatas.CheckConfig
		userToPolociesElevated []UserToPoliciesElevate
		testName               string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if  users can elevate rights",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				userToPolociesElevated: []UserToPoliciesElevate{
					{
						UserName: "test",
						Policies: [][]string{},
					},
				},
				testName: "AWS_IAM_003",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfUserCanElevateRights(tt.args.checkConfig, tt.args.userToPolociesElevated, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfUserCanElevateRights() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfUserCanElevateRightsFAIL(t *testing.T) {
	type args struct {
		checkConfig            yatas.CheckConfig
		userToPolociesElevated []UserToPoliciesElevate
		testName               string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if  users can elevate rights",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				userToPolociesElevated: []UserToPoliciesElevate{
					{
						UserName: "test",
						Policies: [][]string{
							{"test"},
						},
					},
				},
				testName: "AWS_IAM_003",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfUserCanElevateRights(tt.args.checkConfig, tt.args.userToPolociesElevated, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfUserCanElevateRights() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
