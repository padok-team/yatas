package commons

// Wrapper struct that holds all the results of the checks for each account
type Tests struct {
	Account string  `yaml:"account"` // Account name
	Checks  []Check `yaml:"checks"`  // Checks
}
