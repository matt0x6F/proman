package dto

type Editor struct {
	Icon    string `json:"icon,omitempty"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Default bool   `json:"default,omitempty"`
}
