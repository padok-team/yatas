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
