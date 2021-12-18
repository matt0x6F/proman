package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/user"

	"github.com/mattouille/proman/dto"

	"github.com/mitchellh/mapstructure"
	"go.etcd.io/bbolt"
)

var (
	projectBucket = []byte("projects")
	editorBucket  = []byte("editors")

	ErrNoRecords = errors.New("no records found")
)

const DefaultDBPermissions = 0o655

var db *DB

// Service returns an instance of the DB service
func Service() *DB {
	return db
}

// New starts the DB service
func New() error {
	usr, _ := user.Current()

	conn, err := bbolt.Open(usr.HomeDir+"/.config/proman/store.db", DefaultDBPermissions, nil)
	if err != nil {
		return err
	}

	svc := new(DB)
	svc.db = conn
	db = svc

	return db.migrate()
}

// DB is the database service
type DB struct {
	db *bbolt.DB
}

func (d *DB) migrate() error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(projectBucket)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(editorBucket)
		if err != nil {
			return err
		}

		return nil
	})
}

// DeleteProject deletes a project by path.
func (d *DB) DeleteProject(path string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(projectBucket).Delete([]byte(path))
	})
}

// UpsertProject updates or creates a new project. At minimum the "path" key must be provided by proj
func (d *DB) UpsertProject(input map[string]interface{}) error {
	path, ok := input["path"]
	if !ok {
		return fmt.Errorf("path is required to upsert")
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		data := tx.Bucket(projectBucket).Get([]byte(path.(string)))
		// project doesn't exist in the DB
		if len(data) == 0 {
			buff := new(bytes.Buffer)

			err := json.NewEncoder(buff).Encode(input)
			if err != nil {
				return fmt.Errorf("unable to encode project: %w", err)
			}

			return tx.Bucket(projectBucket).Put([]byte(path.(string)), buff.Bytes())
		}

		// decode the project for updating
		var tmp dto.Project

		err := json.NewDecoder(bytes.NewReader(data)).Decode(&tmp)
		if err != nil {
			return fmt.Errorf("unable to decode project: %w", err)
		}

		// update the project
		err = mapstructure.Decode(input, &tmp)
		if err != nil {
			return fmt.Errorf("unable to decode input: %w", err)
		}

		// encode the updated project
		buff := new(bytes.Buffer)

		err = json.NewEncoder(buff).Encode(tmp)
		if err != nil {
			return fmt.Errorf("unable to encode project: %w", err)
		}

		return tx.Bucket(projectBucket).Put([]byte(path.(string)), buff.Bytes())
	})
}

// GetAllProjects fetches all projects from the projects table
func (d *DB) GetAllProjects() ([]dto.Project, error) {
	var projects []dto.Project

	err := d.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(projectBucket).ForEach(func(k, v []byte) error {
			var tmp dto.Project

			err := json.NewDecoder(bytes.NewReader(v)).Decode(&tmp)
			if err != nil {
				return fmt.Errorf("error while decoding %s: %w", k, err)
			}

			projects = append(projects, tmp)

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, ErrNoRecords
	}

	return projects, nil
}

// GetProjectByPath fetches a project by the directory name relative to the project directory
func (d *DB) GetProjectByPath(directory string) (dto.Project, error) {
	var project dto.Project

	err := d.db.View(func(tx *bbolt.Tx) error {
		p := tx.Bucket(projectBucket).Get([]byte(directory))

		// error checking
		if len(p) == 0 {
			return fmt.Errorf("no project found")
		}

		return json.NewDecoder(bytes.NewReader(p)).Decode(&project)
	})
	if err != nil {
		return dto.Project{}, err
	}

	return project, nil
}

// GetEditors fetches the full list of editors
func (d *DB) GetEditors() ([]dto.Editor, error) {
	var editors []dto.Editor

	err := d.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(editorBucket).ForEach(func(k, v []byte) error {
			var tmp dto.Editor

			err := json.NewDecoder(bytes.NewReader(v)).Decode(&tmp)
			if err != nil {
				return fmt.Errorf("error while decoding %s: %w", k, err)
			}

			editors = append(editors, tmp)

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	if len(editors) == 0 {
		return nil, ErrNoRecords
	}

	return editors, nil
}

// UpsertEditor updates or creates a new editor. Path and name are required in order to avoid broken config.
func (d *DB) UpsertEditor(input map[string]interface{}) error {
	_, ok := input["path"]
	if !ok {
		return fmt.Errorf("path is required to upsert")
	}

	name, ok := input["name"]
	if !ok {
		return fmt.Errorf("name is required to upsert")
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		data := tx.Bucket(editorBucket).Get([]byte(name.(string)))
		// editor doesn't exist in the DB
		if len(data) == 0 {
			buff := new(bytes.Buffer)

			err := json.NewEncoder(buff).Encode(input)
			if err != nil {
				return fmt.Errorf("unable to encode editor: %w", err)
			}

			return tx.Bucket(editorBucket).Put([]byte(name.(string)), buff.Bytes())
		}

		// decode the editor for updating
		var tmp dto.Editor

		err := json.NewDecoder(bytes.NewReader(data)).Decode(&tmp)
		if err != nil {
			return fmt.Errorf("unable to decode editor: %w", err)
		}

		// update the editor
		err = mapstructure.Decode(input, &tmp)
		if err != nil {
			return fmt.Errorf("unable to decode input: %w", err)
		}

		// encode the updated editor
		buff := new(bytes.Buffer)

		err = json.NewEncoder(buff).Encode(tmp)
		if err != nil {
			return fmt.Errorf("unable to encode editor: %w", err)
		}

		return tx.Bucket(editorBucket).Put([]byte(name.(string)), buff.Bytes())
	})
}

// DeleteEditor deletes an editor by name
func (d *DB) DeleteEditor(name string) error {
	return d.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(editorBucket).Delete([]byte(name))
	})
}
