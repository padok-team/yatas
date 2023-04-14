package commons

import (
	"sync"
)

// Check test if a wrapper around a check that allows to verify if the check is included or excluded and add some custom logic.
// It allows for a simpler integration of new tests without bloating the code.
func CheckTest[A, B, C any](wg *sync.WaitGroup, config *Config, id string, test func(A, B, C)) func(A, B, C) {
	if !config.CheckExclude(id) && config.CheckInclude(id) {
		wg.Add(1)

		return test
	} else {
		return func(A, B, C) {}
	}

}

// Check Macro test is a wrapper around a category that runs all the checks in the category.
// It allows for a simpler integration of new categories without bloating the code.
func CheckMacroTest[A, B, C, D any](wg *sync.WaitGroup, config *Config, test func(A, B, C, D)) func(A, B, C, D) {
	wg.Add(1)

	return test
}

type CheckFunc func(interface{}) Result

type Resource interface {
	GetID() string
}

type CheckDefinition struct {
	Title          string
	Description    string
	Categories     []string
	ConditionFn    func(Resource) bool
	SuccessMessage string
	FailureMessage string
}

func CheckResources(checkConfig CheckConfig, resources []Resource, checkDefinitions []CheckDefinition) {
	for _, checkDefinition := range checkDefinitions {
		if !checkConfig.ConfigYatas.CheckExclude(checkDefinition.Title) && checkConfig.ConfigYatas.CheckInclude(checkDefinition.Title) {
			check := createCheck(checkDefinition)
			for _, resource := range resources {
				result := checkResource(resource, checkDefinition.ConditionFn, checkDefinition.SuccessMessage, checkDefinition.FailureMessage)
				check.AddResult(result)
			}
			checkConfig.Queue <- check
		}
	}
}

func AddChecks(checkConfig *CheckConfig, resources ...[]CheckDefinition) {
	//Print the resources to check
	totalCount := 0
	for _, resourceSlice := range resources {
		totalCount += len(resourceSlice)
	}
	checkConfig.Wg.Add(totalCount)
}

func createCheck(checkDefinition CheckDefinition) Check {
	var check Check
	check.InitCheck(checkDefinition.Description, checkDefinition.Description, checkDefinition.Title, checkDefinition.Categories)
	return check
}

func checkResource(resource Resource, conditionFn func(Resource) bool, successMessage, failureMessage string) Result {
	if conditionFn(resource) {
		message := successMessage + " - Resource " + resource.GetID()
		return Result{Status: "OK", Message: message, ResourceID: resource.GetID()}
	} else {
		message := failureMessage + " - Resource " + resource.GetID()
		return Result{Status: "FAIL", Message: message, ResourceID: resource.GetID()}
	}
}
