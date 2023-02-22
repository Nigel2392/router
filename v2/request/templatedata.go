package request

import (
	"html/template"
)

type CSRFToken struct {
	Token string
}

func NewCSRFToken(token string) *CSRFToken {
	return &CSRFToken{Token: token}
}

func (ct *CSRFToken) String() string {
	return ct.Token
}

func (ct *CSRFToken) Input() template.HTML {
	return template.HTML("<input type=\"hidden\" name=\"csrf_token\" value=\"" + ct.Token + "\" />")
}

type TemplateData struct {
	Data      map[string]any
	Messages  []Message
	CSRFToken CSRFToken
	User      interface{}
	Next      string
}

func NewTemplateData() *TemplateData {
	return &TemplateData{Data: make(map[string]any), Messages: make([]Message, 0)}
}

func (td *TemplateData) AddMessage(messageType, message string) {
	td.Messages = append(td.Messages, Message{Type: messageType, Text: message})
}

func (td *TemplateData) Set(key string, value interface{}) {
	td.Data[key] = value
}

func (td *TemplateData) Get(key string) interface{} {
	return td.Data[key]
}

func (td *TemplateData) Has(key string) bool {
	_, ok := td.Data[key]
	return ok
}

func (td *TemplateData) Delete(key string) {
	delete(td.Data, key)
}

type Message struct {
	Type string
	Text string
}
