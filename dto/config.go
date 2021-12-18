package dto

// ConfigSchema represents the config keys and values
type ConfigSchema struct {
	ProjectDirectory string `mapstructure:"project_directory" json:"project_directory"`
}
