package dto

// Project is a project managed in the project directory.
type Project struct {
	// Path is the path from the root Project Directory
	Path string `json:"path" mapstructure:"path"`
	// OpenWith is a binary on the system we can use to open the project
	OpenWith string `json:"open_with" mapstructure:"open_with"`
	// Hide means to intentionally hide the project on the main project list
	Hide bool `json:"hide" mapstructure:"hide"`
	// Remotes are the git remotes
	Remotes []string `json:"remotes,omitempty" mapstructure:"remotes"`
	// RepositoryURL is the url to the repository
	RepositoryURLs []string `json:"repository_urls,omitempty" mapstructure:"repository_urls"`
}
