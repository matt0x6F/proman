package main

import (
	"github.com/mattouille/proman/dto"
	"github.com/mattouille/proman/service/config"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

func NewConfig() *Config {
	return new(Config)
}

type Config struct {
	config  *config.Config
	runtime *wails.Runtime
	log     *logger.CustomLogger
}

func (c *Config) WailsInit(runtime *wails.Runtime) error {
	c.runtime = runtime
	c.config = config.Service()
	c.log = c.runtime.Log.New("config")

	c.registerEvents()

	return nil
}

// Registers events which can be called via the wails runtime
func (c *Config) registerEvents() {
	// Updates config based on an event
	c.runtime.Events.On("config.update", func(optionalData ...interface{}) {
		if len(optionalData) == 0 {
			c.log.Error("Frontend attempted to update config but config was blank")

			return
		}

		if optionalData[0] == nil {
			c.log.Error("Frontend attempted to update config but config was blank")

			return
		}

		data := optionalData[0].(map[string]interface{})

		c.log.Debugf("Updating config", data)

		err := c.Update(data)
		if err != nil {
			c.log.Errorf("Error while updating config: ", err)
		}
	})

	c.runtime.Events.On("config.select_project_directory", func(optionalData ...interface{}) {
		dir := c.runtime.Dialog.SelectDirectory()

		// blank return means that the user hit cancel
		if dir != "" {
			err := c.Update(map[string]interface{}{
				"project_directory": dir,
			})

			c.runtime.Events.Emit("config.set_project_directory", dir, err)
		}
	})
}

// Get unmarshals config using mapstructure
func (c *Config) Get() (dto.ConfigSchema, error) {
	return c.config.Unmarshal()
}

// Update reads values from a map[string]interface and writes them to the config file
func (c *Config) Update(data map[string]interface{}) error {
	err := c.config.MergeConfigMap(data)
	if err != nil {
		c.log.Errorf("Error merging config from frontend", err)

		return err
	}

	c.log.Debugf("Updating config", data)

	return c.save()
}

func (c *Config) save() error {
	return c.config.WriteConfig()
}
