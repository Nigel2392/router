package extensions

import "github.com/Nigel2392/router/v3/request"

// Simple extension struct.
// This is used to render the extension into the base template.
// Avoids having to create a new struct for each extension.
type Simple struct {
	// The name of the extension.
	// This is used to uniquely identify the extension.
	ExtensionName string

	// The file name for the extension.
	// This is the name of the template to render.
	FileName string

	// The callback that is called when the extension is rendered.
	// This is used to get the template data and the template name.
	// The template name is the name of the template to render.
	// This is not the file name, but the name of the template.
	// IE: {{template "base" .}}
	Callback func(*request.Request) map[string]any
}

// Returns the name of the extension.
func (s *Simple) Name() string {
	return s.ExtensionName
}

// Returns the file name of the extension.
// This is the name of the template to render.
// The template will be fetched from a template.Manager.
func (s *Simple) Filename() string {
	return s.FileName
}

// Returns the template data for the extension.
func (s *Simple) View(r *request.Request) map[string]any {
	return s.Callback(r)
}
