package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/mattouille/proman/dto"
	"github.com/spf13/viper"
)

const (
	ConfigPath = "/.config/proman"
	ConfigName = "/config.toml"
)

var (
	c *Config
)

func init() {
	c = New()
}

// New creates a new config service
func New() *Config {
	c = new(Config)
	c.viper = viper.New()

	return c
}

// Service returns an instance of the config service.
func Service() *Config {
	return c
}

// Config is the Config Service
type Config struct {
	viper *viper.Viper
}

// ReadInConfig reads configuration from a specified location
func ReadInConfig() error { return c.ReadInConfig() }

func (c *Config) ReadInConfig() error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("unable to determine current user: %w", err)
	}

	home := usr.HomeDir

	path := home + ConfigPath
	cfg := path + ConfigName

	// make sure the config directory exists
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0o775) //nolint:gomnd
		if err != nil {
			return fmt.Errorf("unable to create proman directory: %w", err)
		}
	}

	// make sure the config file exists
	_, err = os.Stat(cfg)
	if os.IsNotExist(err) {
		file, err := os.Create(cfg)
		if err != nil {
			return fmt.Errorf("unable to create config file: %w", err)
		}

		err = file.Close()
		if err != nil {
			return fmt.Errorf("unable to close config file: %w", err)
		}
	}

	c.viper.SetConfigFile(cfg)

	err = c.viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("unable to read in config: %w", err)
	}

	return nil
}

// Unmarshal unmarshals config using mapstructure
func Unmarshal() (dto.ConfigSchema, error) { return c.Unmarshal() }

func (c *Config) Unmarshal() (dto.ConfigSchema, error) {
	cfg := dto.ConfigSchema{}

	err := c.viper.Unmarshal(&cfg)
	if err != nil {
		return dto.ConfigSchema{}, err
	}

	return cfg, err
}

// MergeConfigMap reads values from a map[string]interface and writes them to the config file
func MergeConfigMap(data map[string]interface{}) error { return c.MergeConfigMap(data) }

func (c *Config) MergeConfigMap(data map[string]interface{}) error {
	return c.viper.MergeConfigMap(data)
}

// WriteConfig writes the configuration back to disk
func WriteConfig() error { return c.WriteConfig() }

func (c *Config) WriteConfig() error {
	return c.viper.WriteConfig()
}
