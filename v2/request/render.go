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

	r.Data.Next = r.Next()

	// Get the messages from the session
	if r.Session != nil {
		if r.Data.Messages == nil {
			r.Data.Messages = make([]Message, 0)
		}
		if messages, ok := r.Session.Get("messages").([]Message); ok {
			r.Data.Messages = append(r.Data.Messages, messages...)
			r.Session.Delete("messages")
		}
	}

	// Add default data
	if DEFAULT_DATA_FUNC != nil {
		DEFAULT_DATA_FUNC(r)
	}
	// Render template
	return t.ExecuteTemplate(r.Response, name, r.Data)
}
