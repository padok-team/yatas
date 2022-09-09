package commons

type Ignore struct {
	ID     string   `yaml:"id"`
	Regex  bool     `yaml:"regex"`
	Values []string `yaml:"values"`
}
