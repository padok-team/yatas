package guardduty

import (
	"sync"
	"testing"

	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfGuarddutyEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		testName    string
		detectors   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfGuarddutyEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				testName:  "TestCheckIfGuarddutyEnabled",
				detectors: []string{"detector1", "detector2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfGuarddutyEnabled(tt.args.checkConfig, tt.args.testName, tt.args.detectors)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckifGuarddutyEnabled() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfGuarddutyEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		testName    string
		detectors   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfGuarddutyEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				testName:  "TestCheckIfGuarddutyEnabled",
				detectors: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfGuarddutyEnabled(tt.args.checkConfig, tt.args.testName, tt.args.detectors)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckifGuarddutyEnabled() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
