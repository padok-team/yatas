package commons

import (
	"sync"
)

// CheckConfig is a struct that contains all the information needed to run a check.
type CheckConfig struct {
	Wg          *sync.WaitGroup // Wait group to wait for all the checks to be done
	Queue       chan Check      // Queue to add the results to
	ConfigYatas *Config         // Yatas config
}

// Init the check config struct. Particularly useful in the categories. It allows to pass the config to the checks and allows
// them to be run in parallel by adding the results to the queue.
func (c *CheckConfig) Init(config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}
