package request

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"html/template"
)

type TemplateRequest struct {
	url  func(string, ...interface{}) string
	User User
	Next string
}

func (tr *TemplateRequest) URL(path string, args ...interface{}) string {
	return tr.url(path, args...)
}

func init() {
	gob.Register(Messages{})
}

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
	Messages  Messages
	CSRFToken *CSRFToken
	Request   *TemplateRequest
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Data:      make(map[string]any),
		Messages:  make(Messages, 0),
		CSRFToken: nil,
		Request:   &TemplateRequest{},
	}
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

type Message struct {
	Type string
	Text string
}

type Messages []Message

func (m Messages) Encode() string {
	var buf = bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	enc.Encode(m)
	// Encode with base64
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (m *Messages) Decode(data string) error {
	// Decode with base64
	var buf, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	return dec.Decode(m)
}
