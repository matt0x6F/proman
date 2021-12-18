package main

import (
	"github.com/mattouille/proman/path"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

// Validate does various forms of validation
type Validate struct {
	runtime *wails.Runtime
	log     *logger.CustomLogger
}

func NewValidator() *Validate {
	return new(Validate)
}

func (v *Validate) WailsInit(runtime *wails.Runtime) error {
	v.runtime = runtime
	v.log = v.runtime.Log.New("validation")

	v.registerEvents()

	return nil
}

// Registers events which can be called via the wails runtime
func (v *Validate) registerEvents() {
	// validate.config validates an entire configuration payload at once and emits an event with the errors payload
	v.runtime.Events.On("validate.config", func(optionalData ...interface{}) {
		data := optionalData[0].(map[string]interface{})
		errors := v.Configuration(data)

		v.runtime.Events.Emit("validate.config.completed", errors)
	})
}

// Configuration takes in a map[string]interface{} configuration and returns an error map where keys are property names
func (v *Validate) Configuration(cfg map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	ok, err := v.ProjectDir(cfg["project_directory"].(string))
	if !ok {
		errors["project_directory"] = err.Error()
	}

	return errors
}

// ProjectDir validates a project directory
func (v *Validate) ProjectDir(projectDir string) (bool, error) {
	abs, err := path.ExpandAndValidate(projectDir)
	if err != nil {
		v.log.DebugFields("Validation failed on project directory", logger.Fields{"path": projectDir, "error": err})

		return false, err
	}

	v.log.DebugFields("Target directory exists", logger.Fields{"path": abs})

	return true, nil
}
