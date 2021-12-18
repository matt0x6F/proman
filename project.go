package main

import (
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/mattouille/proman/dto"
	"github.com/mattouille/proman/path"
	"github.com/mattouille/proman/service/config"
	"github.com/mattouille/proman/service/database"
	"github.com/wailsapp/wails"
	"github.com/wailsapp/wails/lib/logger"

	"github.com/go-git/go-git/v5"
)

func NewProjects() *Projects {
	return new(Projects)
}

// Projects is the Projects frontend service.
type Projects struct {
	runtime  *wails.Runtime
	log      *logger.CustomLogger
	projects []dto.Project
}

func (p *Projects) WailsInit(runtime *wails.Runtime) error {
	p.runtime = runtime
	p.log = p.runtime.Log.New("project")

	cfg, err := config.Unmarshal()
	if err != nil {
		return err
	}

	paths, err := p.loadProjectsFromDisk(cfg.ProjectDirectory)
	if err != nil {
		return err
	}

	p.projects, err = p.syncProjectMetadata(paths)
	if err != nil {
		return err
	}

	p.registerEvents()

	return nil
}

// Registers events which can be called via the wails runtime
func (p *Projects) registerEvents() {
	p.runtime.Events.On("OpenURL", func(optionalData ...interface{}) {
		if len(optionalData) == 0 {
			p.log.Error("Frontend attempted to open a URL but the URL was blank")

			return
		}

		if optionalData[0] == nil {
			p.log.Error("Frontend attempted to open a URL but the URL was blank")

			return
		}

		url := optionalData[0].(string)

		p.log.DebugFields("Opening URL", logger.Fields{"url": url})

		err := p.runtime.Browser.OpenURL(url)
		if err != nil {
			p.log.ErrorFields("Error while opening URL", logger.Fields{"error": err})
		}
	})

	p.runtime.Events.On("OpenProject", func(optionalData ...interface{}) {
		if len(optionalData) == 0 {
			p.log.Error("Frontend attempted to open a URL but the URL was blank")

			return
		}

		if optionalData[0] == nil {
			p.log.Error("Frontend attempted to open a URL but the URL was blank")

			return
		}

		path := optionalData[0]

		p.log.DebugFields("Opening project", logger.Fields{"path": path})
	})
}

// Returns a slice of paths known to be project directories in the project directory path
func (p *Projects) loadProjectsFromDisk(projectDir string) ([]string, error) {
	var projects []string

	abs, err := path.ExpandAndValidate(projectDir)
	if err != nil {
		p.log.ErrorFields("Failed loading project directory", logger.Fields{"path": projectDir, "error": err})

		return nil, err
	}

	files, err := ioutil.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	// create or update project
	for _, f := range files {
		if f.IsDir() {
			projects = append(projects, f.Name())

			p.log.DebugFields("Searching repository git path", logger.Fields{"path": abs + "/" + f.Name()})

			// todo: in the future, more vcs providers could be supported. it might be worth splitting vcs detection
			// into it's own
			repo, err := git.PlainOpen(abs + "/" + f.Name())
			if err != nil {
				if !errors.Is(err, git.ErrRepositoryNotExists) {
					return nil, err
				}

				p.log.Debug("Git repository not detected")
			}

			var (
				remotes []string
				urls    []string
			)

			if repo != nil {
				p.log.Debug("Git repository detected")

				rmts, err := repo.Remotes()
				if err != nil {
					return nil, err
				}

				for _, remote := range rmts {
					remotes = append(remotes, remote.Config().URLs...)
					urls = append(urls, p.ParseRepositoryURL(remote.Config().URLs)...)
				}
			}

			p.log.DebugFields("Found project", logger.Fields{"name": f.Name(), "remotes": remotes, "urls": urls})

			err = database.Service().UpsertProject(map[string]interface{}{"path": f.Name(), "remotes": remotes, "repository_urls": urls})
			if err != nil {
				p.log.ErrorFields("Error while upserting project", logger.Fields{"error": err})
			}
		}
	}

	return projects, nil
}

// Retrieves a list of all active projects from the database, matches them by path, and sets any defaults on new metadata.
// paths should be relative to the project directory.
func (p *Projects) syncProjectMetadata(paths []string) ([]dto.Project, error) {
	projects, err := database.Service().GetAllProjects()
	if err != nil && !errors.Is(err, database.ErrNoRecords) {
		return nil, err
	}

	var (
		final []dto.Project
	)

	// add projects from the db that are in the current path list
	for _, projectPath := range paths {
		for i := 0; i < len(projects); i++ {
			if projects[i].Path == projectPath {
				final = append(final, projects[i])
				// shift the element to the back
				projects[i] = projects[len(projects)-1]

				break
			}
		}
	}

	// delete projects which no longer exist
	del := len(projects) - len(final)
	if del < 1 {
		del = 1
	}

	projects = projects[del-1:]

	for i := 0; i < len(projects); i++ {
		err := database.Service().DeleteProject(projects[i].Path)
		if err != nil {
			p.log.ErrorFields("Unable to delete project", logger.Fields{"path": projects[i].Path})
		}
	}

	return final, nil
}

// GetAll fetches all projects from the database
func (p *Projects) GetAll(refresh bool) ([]dto.Project, error) {
	if refresh {
		cfg, err := config.Unmarshal()
		if err != nil {
			return nil, err
		}

		paths, err := p.loadProjectsFromDisk(cfg.ProjectDirectory)
		if err != nil {
			return nil, err
		}

		p.projects, err = p.syncProjectMetadata(paths)
		if err != nil {
			return nil, err
		}
	}

	return p.projects, nil
}

// ParseRepositoryURL parses remote urls and returns the expected repository urls. If one cannot be determined it returns
// nil.
//
// Remotes can be in the form of http(s) or ssh.
func (p *Projects) ParseRepositoryURL(urls []string) []string {
	var repos []string

	for _, url := range urls {
		p.log.DebugFields("Parsing URL", logger.Fields{"url": url})

		// todo: support more vcs providers
		// if the go compiler is any good indicator, building a regex for this could be complicated but worthwhile.
		// it might even be possible for someone to add their own regex and logo for specific or custom providers.
		//
		// an alternative would be to allow setting the regex that an individual project uses and fetching their
		// favicon automatically.

		// parses ssh, http, and https for github.com using git
		re := regexp.MustCompile(`^.+(@|://)(?P<service>.+)[:/](?P<user>.+)/(?P<repo>.+)\.(?P<vcs>.+)$`)
		if !re.MatchString(url) {
			p.log.DebugFields("Failed to parse URL", logger.Fields{"url": url})

			continue
		}

		p.log.DebugFields("Matched URL", logger.Fields{"url": url})

		matches := re.FindStringSubmatch(url)
		svc := re.SubexpIndex("service")

		if matches[svc] == "github.com" {
			usr := re.SubexpIndex("user")
			repo := re.SubexpIndex("repo")
			repos = append(repos, "https://"+matches[svc]+"/"+matches[usr]+"/"+matches[repo]+"/")

			continue
		}
	}

	return repos
}
