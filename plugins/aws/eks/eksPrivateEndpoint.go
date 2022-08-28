package eks

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/internal/logger"
	"github.com/stangirard/yatas/internal/yatas"
	"golang.org/x/exp/slices"
)

func CheckIfEksEndpointPrivate(checkConfig yatas.CheckConfig, clusters []types.Cluster, testName string) {
	logger.Info(fmt.Sprint("Running ", testName))
	var check yatas.Check
	check.InitCheck("EKS clusters have private endpoint or strict public access", "Check if EKS clusters have private endpoint", testName)
	for _, cluster := range clusters {
		if cluster.ResourcesVpcConfig != nil {
			if cluster.ResourcesVpcConfig.EndpointPublicAccess {
				if ok := slices.Contains(cluster.ResourcesVpcConfig.PublicAccessCidrs, "0.0.0.0/0"); !ok {
					Message := "EKS cluster " + *cluster.Name + " has private endpoint"
					result := yatas.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				} else {
					Message := "EKS cluster " + *cluster.Name + " has public endpoint"
					result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
					check.AddResult(result)
				}
			} else {
				Message := "EKS cluster " + *cluster.Name + " has private endpoint"
				result := yatas.Result{Status: "OK", Message: Message, ResourceID: *cluster.Name}
				check.AddResult(result)
			}
		} else {
			Message := "Private endpoint is not enabled for cluster " + *cluster.Name
			result := yatas.Result{Status: "FAIL", Message: Message, ResourceID: *cluster.Name}
			check.AddResult(result)
		}
	}
	checkConfig.Queue <- check
}
