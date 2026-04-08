package utils

import (
	"bytes"
	"html/template"
)

func RenderTemplate(path string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(path);
	if err != nil {
		return "" , err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf,data)
	if err != nil {
		return "" , nil
	}
	return buf.String() , nil
}