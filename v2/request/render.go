package request

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nigel2392/router/v2/templates"
)

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
		var messages = r.Session.Get(MESSAGE_COOKIE_NAME)
		if messages != nil {
			fmt.Println("Messages:", messages)
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
