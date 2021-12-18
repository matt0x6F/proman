package main

import (
	"errors"

	"github.com/mattouille/proman/dto"
	"github.com/mattouille/proman/service/database"

	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"
)

func NewEditorConfig() *EditorConfig {
	return new(EditorConfig)
}

type EditorConfig struct {
	Editors []dto.Editor `json:"editors"`
	runtime *wails.Runtime
	log     *logger.CustomLogger
	db      *database.DB
}

func (c *EditorConfig) WailsInit(runtime *wails.Runtime) error {
	c.runtime = runtime
	c.log = c.runtime.Log.New("config")
	c.db = database.Service()

	// preload the editors
	_, err := c.GetAll(true)
	if err != nil && !errors.Is(err, database.ErrNoRecords) {
		c.log.ErrorFields("Error while retrieving editors", logger.Fields{"error": err})
		return err
	}

	c.registerEvents()

	return nil
}

// Registers events which can be called via the wails runtime
func (c *EditorConfig) registerEvents() {
	c.runtime.Events.On("editor.upsert", func(optionalData ...interface{}) {
		if len(optionalData) == 0 {
			c.log.Error("Frontend attempted to update an editor but the editor was blank")

			return
		}

		if optionalData[0] == nil {
			c.log.Error("Frontend attempted to update an editor but the editor was blank")

			return
		}

		data := optionalData[0].(map[string]interface{})

		err := c.UpsertEditor(data)
		if err != nil {
			c.log.ErrorFields("Error while upserting editor", logger.Fields{"error": err})
		}
	})

	c.runtime.Events.On("editor.remove", func(optionalData ...interface{}) {
		if len(optionalData) == 0 {
			c.log.Error("Frontend attempted to remove an editor but the editor was blank")

			return
		}

		if optionalData[0] == nil {
			c.log.Error("Frontend attempted to remove an editor but the editor was blank")

			return
		}

		data := optionalData[0].(string)

		err := c.RemoveEditor(data)
		if err != nil {
			c.log.ErrorFields("Error while removing editor", logger.Fields{"error": err})
		}
	})
}

func (c *EditorConfig) GetAll(refresh bool) ([]dto.Editor, error) {
	if refresh {
		editors, err := c.db.GetEditors()
		if err != nil {
			return nil, err
		}

		c.Editors = editors
	}

	return c.Editors, nil
}

func (c *EditorConfig) UpsertEditor(data map[string]interface{}) error {
	return c.db.UpsertEditor(data)
}

func (c *EditorConfig) RemoveEditor(name string) error {
	return c.db.DeleteEditor(name)
}
