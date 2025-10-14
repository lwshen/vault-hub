package email

import (
	"bytes"
	"html/template"
)

type TemplateData struct {
	Subject   string
	AppName   string
	UserName  string
	ActionURL string
	TTL       string
}

func renderTemplate(name string, data TemplateData) (string, error) {
	layout, err := template.ParseFS(templatesFS, "templates/layout.html.tmpl")
	if err != nil {
		return "", err
	}
	t, err := layout.ParseFS(templatesFS, "templates/"+name)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
