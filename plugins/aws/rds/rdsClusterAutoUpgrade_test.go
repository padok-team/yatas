package rds

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfClusterAutoUpgradeEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		clusters    []types.DBCluster
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfAutoUpgradeEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				clusters: []types.DBCluster{
					{
						DBClusterIdentifier:     aws.String("test"),
						DBClusterArn:            aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:        true,
						AutoMinorVersionUpgrade: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfClusterAutoUpgradeEnabled(tt.args.checkConfig, tt.args.clusters, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfAutoUpgrade() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func Test_checkIfClusterAutoUpgradeEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		clusters    []types.DBCluster
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test_checkIfAutoUpgradeEnabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				clusters: []types.DBCluster{
					{
						DBClusterIdentifier:     aws.String("test"),
						DBClusterArn:            aws.String("arn:aws:rds:us-east-1:123456789012:db:test"),
						StorageEncrypted:        true,
						AutoMinorVersionUpgrade: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfClusterAutoUpgradeEnabled(tt.args.checkConfig, tt.args.clusters, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					t.Logf("%v", check)
					if check.Status != "FAIL" {
						t.Errorf("CheckIfAutoUpgrade() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
