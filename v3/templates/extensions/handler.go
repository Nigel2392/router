package extensions

import (
	"bytes"
	"html/template"

	"github.com/Nigel2392/router/v3/request"
)

// View returns a view that renders the extension into the base template.
func View(options *Options, ext Extension) func(r *request.Request) {
	return func(r *request.Request) {
		var err error
		var buf bytes.Buffer
		var tdata = ext.View(r)
		var tmpl *template.Template

		switch ext := ext.(type) {
		case ExtensionWithTemplate:
			tmpl = ext.Template(r)
			options.render(&buf, ext, tmpl.Tree.Root.String())
		case ExtensionWithStrings:
			options.render(&buf, ext, ext.String(r))
		default:
			tmpl, err = template.ParseFS(options.ExtensionManager.TEMPLATEFS, ext.Filename())
			if err != nil {
				defaultErr(options, r, err)
				return
			}
			options.render(&buf, ext, tmpl.Tree.Root.String())
		}

		t, err := options.BaseManager.GetFromString(buf.String(), "ext")
		if err != nil {
			defaultErr(options, r, err)
			return
		}

		base, err := options.BaseManager.GetBases(nil)
		if err != nil {
			defaultErr(options, r, err)
			return
		}
		for _, b := range base.Templates() {
			t.AddParseTree(b.Name(), b.Tree)
		}

		t.Funcs(options.BaseManager.DEFAULT_FUNCS)
		t.Funcs(options.ExtensionManager.DEFAULT_FUNCS)

		for k, v := range tdata {
			r.Data.Set(k, v)
		}

		if options.BeforeRender != nil {
			options.BeforeRender(r, t)
		}

		err = t.Execute(r, r.Data)
		if err != nil {
			defaultErr(options, r, err)
			return
		}
	}
}

func defaultErr(o *Options, r *request.Request, err error) {
	if o.OnError != nil {
		o.OnError(r, err)
	} else {
		r.Error(500, err.Error())
	}
}
