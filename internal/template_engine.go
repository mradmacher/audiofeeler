package audiofeeler

import (
	"html/template"
	"net/http"
)

type Template struct {
	tpl *template.Template
}

func (t *Template) Execute(w http.ResponseWriter, r *http.Request, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	templateName := "application"
	if r.Header.Get("Hx-Request") == "true" {
		templateName = "content"
	}

	err := t.tpl.ExecuteTemplate(w, templateName, data)

	return err
}

type TemplateEngine struct {
	templatesPath string
}

func NewTemplateEngine(templatesPath string) TemplateEngine {
	return TemplateEngine{templatesPath}
}

func (t *TemplateEngine) MustParse(names ...string) Template {
	files := []string{t.templatesPath + "/application.gohtml"}
	for _, name := range names {
		files = append(files, t.templatesPath+"/"+name+".gohtml")
	}
	return Template{template.Must(template.ParseFiles(files...))}
}
