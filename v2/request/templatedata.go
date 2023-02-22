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
	CSRFToken *CSRFToken
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
	// td.makeMap()
	td.Data[key] = value
}

func (td *TemplateData) Get(key string) interface{} {
	// td.makeMap()
	return td.Data[key]
}

func (td *TemplateData) Has(key string) bool {
	// td.makeMap()
	_, ok := td.Data[key]
	return ok
}

func (td *TemplateData) Delete(key string) {
	// td.makeMap()
	delete(td.Data, key)
}

//	func (td *TemplateData) makeMap() {
//		if td.Data == nil {
//			td.Data = make(map[string]any)
//		}
//	}

type Message struct {
	Type string
	Text string
}