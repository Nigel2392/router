package request

import "github.com/Nigel2392/router/v2/templates"

var DEFAULT_DATA_FUNC func(r *Request)

// Template configuration must be set before calling this function!
// See the templates package for more information.
func (r *Request) Render(templateName string) error {
	var t, name, err = templates.GetTemplate(templateName)
	if err != nil {
		return err
	}
	if DEFAULT_DATA_FUNC != nil {
		DEFAULT_DATA_FUNC(r)
	}
	// Render template
	return t.ExecuteTemplate(r.Response, name, r.Data)
}
