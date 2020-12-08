package config

/*
define a application
*/

type ApplicationConfiguration struct {
	Services []string `json:"services" yaml:"services"`
}

func LoadApplicationConfig(path string) (*ApplicationConfiguration, error) {
	return &ApplicationConfiguration{Services: []string{}}, nil
}
