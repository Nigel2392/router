package request

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nigel2392/router/v2/templates"
)

var DEFAULT_DATA_FUNC func(r *Request)

var TEMPLATE_MANAGER *templates.Manager

// Template configuration must be set before calling this function!
// See the templates package for more information.
func (r *Request) Render(templateName string) error {
	if TEMPLATE_MANAGER == nil {
		panic("Template manager is nil, please set the template manager before calling Render()")
	}
	var t, name, err = TEMPLATE_MANAGER.Get(templateName)
	if err != nil {
		return err
	}

	r.Data.Next = r.Next()

	// Get the messages from the session
	if r.Session != nil {
		var messages = r.Session.Get(MESSAGE_COOKIE_NAME)
		if messages != nil {
			var messagesCasted, ok = messages.(Messages)
			if !ok {
				return fmt.Errorf("Messages in session are not of type Messages, but %T", messages)
			}
			r.Data.Messages = append(r.Data.Messages, messagesCasted...)
			r.Session.Delete(MESSAGE_COOKIE_NAME)
		}
	} else {
		if cookie, err := r.Request.Cookie(MESSAGE_COOKIE_NAME); err == nil {
			(&r.Data.Messages).Decode(cookie.Value)
			// Delete the cookie.
			http.SetCookie(r.Response, &http.Cookie{
				Name:    MESSAGE_COOKIE_NAME,
				Value:   "",
				Expires: time.Now().Add(-time.Hour),
			})
		}
	}

	// Add default data
	if DEFAULT_DATA_FUNC != nil {
		DEFAULT_DATA_FUNC(r)
	}

	// Render template
	return t.ExecuteTemplate(r.Response, name, r.Data)
}
