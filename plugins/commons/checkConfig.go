package commons

import (
	"sync"
)

// Init the check config struct. Particularly useful in the categories. It allows to pass the config to the checks and allows
// them to be run in parallel by adding the results to the queue.
func (c *CheckConfig) Init(config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}
