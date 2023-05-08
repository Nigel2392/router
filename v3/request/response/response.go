package response

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/templates"
)

var DEFAULT_DATA_FUNC func(r *request.Request)

var TEMPLATE_MANAGER *templates.Manager

// Template configuration must be set before calling this function!
//
// See the templates package for more information.
func Render(r *request.Request, templateName string) error {
	if TEMPLATE_MANAGER == nil {
		panic("Template manager is nil, please set the template manager before calling Render()")
	}
	var t, name, err = TEMPLATE_MANAGER.Get(templateName)
	if err != nil {
		return err
	}

	return Template(r, t, name)
}

// Render a string as a template
func String(r *request.Request, templateString string) error {
	var t = template.New("string")
	t, err := t.Parse(templateString)
	if err != nil {
		return err
	}
	// Render template
	return Template(r, t, "string")
}

// Render a template with the given name
func Template(r *request.Request, t *template.Template, name string) error {
	var err = addDefaultData(r)
	if err != nil {
		return err
	}

	// Add default data
	if DEFAULT_DATA_FUNC != nil {
		DEFAULT_DATA_FUNC(r)
	}
	// Render template
	return t.ExecuteTemplate(r.Response, name, r.Data)
}

// Add default data to the request
func addDefaultData(r *request.Request) error {
	if r.Data == nil {
		r.Data = &request.TemplateData{}
	}
	if r.Data.Request == nil {
		r.Data.Request = &request.TemplateRequest{}
	}
	r.Data.Request.Next = r.Next()
	r.Data.Request.User = r.User
	// Get the messages from the session
	if r.Session != nil {
		var messages = r.Session.Get(request.MESSAGE_COOKIE_NAME)
		if messages != nil {
			var messagesCasted, ok = messages.(request.Messages)
			if !ok {
				//lint:ignore ST1005 This is a log message, not an error message.
				return fmt.Errorf("Messages in session are not of type Messages, but %T", messages)
			}
			r.Data.Messages = append(r.Data.Messages, messagesCasted...)
			r.Session.Delete(request.MESSAGE_COOKIE_NAME)
		}
	} else {
		if cookie, err := r.Request.Cookie(request.MESSAGE_COOKIE_NAME); err == nil {
			(&r.Data.Messages).Decode(cookie.Value)
			// Delete the cookie.
			http.SetCookie(r.Response, &http.Cookie{
				Name:    request.MESSAGE_COOKIE_NAME,
				Value:   "",
				Expires: time.Now().Add(-time.Hour),
			})
		}
	}
	return nil
}
