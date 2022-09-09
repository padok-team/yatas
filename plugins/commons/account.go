package commons

// AWS Account struct
type AWS_Account struct {
	Name    string `yaml:"name"`    // Name of the account in the reports
	Profile string `yaml:"profile"` // Profile to use
	SSO     bool   `yaml:"sso"`     // Use SSO
	Region  string `yaml:"region"`  // Region to use
}
