package extensions

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/templates"
)

// ExtensionWithTemplate is an extension that has a template.
// This is an addition to the default extensions, where you can specify your own template.
type ExtensionWithTemplate interface {
	Extension
	Template(*request.Request) *template.Template
}

// ExtensionWithStrings is an extension that exists of a string.
// This is an addition to the default extensions, where you can specify your own string as a template.
type ExtensionWithStrings interface {
	Extension
	String(*request.Request) string
}

// Template extensions
// These are extensions that are rendered into the base template.
// This is useful, if you want to allow people from other packages
// to extend base templates of for example; an admin panel.
type Extension interface {

	// The name of the extension.
	// This is used to uniquely identify the extension.
	Name() string

	// The file name for the extension.
	// This is the name of the template to render.
	// The template will be fetched from a template.Manager.
	Filename() string

	// Extra data for the extension when it is rendered.
	View(*request.Request) map[string]any
}

// Block is a block of text that is rendered before or after the template.
type Block struct {
	// The name of the block.
	// IE: {{define "content"}}, the block name is "content".
	BlockName string
	// The value of the block.
	// IE: {{define "content"}}Hello World{{end}}, the value is "Hello World".
	Value string
}

func (b *Block) write(buf *bytes.Buffer) {
	if b.BlockName != "" {
		buf.WriteString(fmt.Sprintf(`{{define "%s"}}`, b.BlockName))
	}

	buf.WriteString(b.Value)

	if b.BlockName != "" {
		buf.WriteString(`{{end}}`)
	}
}

// Default options for the extension view.
// This is used to render the extension into the base template.
type Options struct {
	// BaseManager is the manager for the base template.
	// It is used to get the base template, and to parse the extension template.
	BaseManager *templates.Manager

	// ExtensionManager is the manager for the extension template.
	ExtensionManager *templates.Manager

	// TemplateName is the name of the template to use as the base.
	// This is not the file name, but the name of the template.
	// IE: {{template "base" .}}
	TemplateName string

	// BlockName is the name of the block to render the template into.
	// IE: {{define "content"}}
	BlockName string

	// Called when an error occurs.
	OnError func(*request.Request, error)

	// Called before the template is rendered.
	BeforeRender func(*request.Request, *template.Template)

	// Custom blocks
	// These are blocks that are rendered before and after the template.
	CSS *Block
	JS  *Block
}

func (o *Options) render(buf *bytes.Buffer, ext Extension, templateString string) {
	buf.WriteString(fmt.Sprintf(`{{template "%s" .}}`, o.TemplateName))
	if o.CSS != nil {
		if o.CSS.BlockName != "" {
			o.CSS.write(buf)
		}
	}
	buf.WriteString(fmt.Sprintf(`{{define "%s"}}`, o.BlockName))
	if o.CSS != nil {
		if o.CSS.BlockName == "" {
			o.CSS.write(buf)
		}
	}
	buf.WriteString(templateString)
	if o.JS != nil {
		if o.JS.BlockName == "" {
			o.JS.write(buf)
		}
	}
	buf.WriteString(`{{end}}`)
	if o.JS != nil {
		if o.JS.BlockName == "" {
			o.JS.write(buf)
		}
	}
}
