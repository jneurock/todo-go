package web

import (
	"fmt"
	"text/template"
)

type Templates struct {
	base      string
	templates map[string]*template.Template
}

func NewTemplates(base string) *Templates {
	return &Templates{
		base:      base,
		templates: map[string]*template.Template{},
	}
}

func (t *Templates) Load(path string) *template.Template {
	fullName := fmt.Sprintf("%s/%s", t.base, path)

	tmpl, ok := t.templates[fullName]

	if !ok {
		tmpl = template.Must(template.ParseGlob(fmt.Sprintf("%s/%s", t.base, path)))
		t.templates[fullName] = tmpl
	}

	return tmpl
}
