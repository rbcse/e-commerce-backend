package templates

import (
	"bytes"
	"html/template"
)

type TemplateRenderer interface {
	Render(path string, msg interface{}) (string, error)
}

type HTMLRenderer struct{}

func (h *HTMLRenderer) Render(path string, msg interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path);
	if err != nil {
		return "" , err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf,msg)
	if err != nil {
		return "" , nil
	}
	return buf.String() , nil
}