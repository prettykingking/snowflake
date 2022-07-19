package config

import (
	"github.com/traefik/paerser/cli"
	"github.com/traefik/paerser/file"
)

// FileLoader loads a configuration from a file.
type FileLoader struct {
	filename   string
	BasePaths  []string
	Extensions []string
}

// GetFilename returns the configuration file if any.
func (f *FileLoader) GetFilename() string {
	return f.filename
}

// Load loads the command's configuration from a file either specified
// with the ConfigFileFlag flag, or from default locations.
func (f *FileLoader) Load(configPath string, config *Configuration) (bool, error) {
	configFile, err := f.loadConfigFiles(configPath, config)
	if err != nil {
		return false, err
	}

	f.filename = configFile

	if configFile == "" {
		return false, nil
	}

	return true, nil
}

// loadConfigFiles tries to decode the given configuration file
// and all default locations for the configuration file.
// It stops as soon as decoding one of them is successful.
func (f *FileLoader) loadConfigFiles(configFile string, element interface{}) (string, error) {
	finder := cli.Finder{
		BasePaths:  []string{"/etc/snowflake", "$HOME/.config/snowflake/snowflake", "./snowflake"},
		Extensions: []string{"toml"},
	}

	filePath, err := finder.Find(configFile)
	if err != nil {
		return "", err
	}

	if len(filePath) == 0 {
		return "", nil
	}

	if err = file.Decode(filePath, element); err != nil {
		return "", err
	}

	return filePath, nil
}
