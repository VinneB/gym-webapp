package middleware

import (
	"bytes"
	"html/template"
	"io"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	t := template.New("main")
	t.Funcs(template.FuncMap{
		"renderTemplate": func(name string, data interface{}) (template.HTML, error) {
			buf := new(bytes.Buffer)
			// Execute the specific template by name into the buffer
			err := t.ExecuteTemplate(buf, name, data)
			if err != nil {
				return "", err
			}
			// Return as html.HTML to prevent further escaping
			return template.HTML(buf.String()), nil
		},
	})
	t, err := t.ParseGlob("pages/template_test.html")
	if err != nil {
		panic(1)
	}
	return &Templates{
		templates: t,
	}
}
