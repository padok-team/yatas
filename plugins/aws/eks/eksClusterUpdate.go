package eks

import (
	"fmt"

	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
)

func CheckIfEKSUpdateAvailable(checkConfig yatas.CheckConfig, clusters []ClusterToUpdate, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("EKS clusters have update available", "Check if EKS update is available", testName)
	for _, cluster := range clusters {
		found := false
		count := 0
		for _, update := range cluster.Updates {
			if update.CreatedAt == nil {
				fmt.Println(update.Type)
				for _, v := range update.Params {
					fmt.Println(v.Type, *v.Value)
				}
				fmt.Println(update.Status)

				found = true
				count++
			}
		}
		if found {
			Message := fmt.Sprintf("%d update(s) available for cluster %s", count, cluster.ClusterName)
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: cluster.ClusterName}
			check.AddResult(result)
		} else {
			Message := "No update available for cluster " + cluster.ClusterName
			result := yatas.Result{Status: "OK", Message: Message, ResourceID: cluster.ClusterName}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
