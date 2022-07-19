package config

// Configuration application configuration
type Configuration struct {
	Logging  *Logging
	Server   *Server
	Settings *Settings
}

// Logging represents loggers configuration
type Logging struct {
	Level string
}

// Server configuration
type Server struct {
	Port uint16
}

// Settings configuration
type Settings struct {
	StartTime string
	MachineId uint16
}

// NewConfiguration returns application default configuration
func NewConfiguration() *Configuration {
	return &Configuration{}
}

func LoadFile(configFilePath string, configuration *Configuration) (bool, error) {
	loader := FileLoader{}
	ok, err := loader.Load(configFilePath, configuration)
	return ok, err
}
