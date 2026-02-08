package audiofeeler

import (
	"html/template"
	"net/http"
)

type TemplateEngine struct {
	templatesPath string
}

func NewTemplateEngine(templatesPath string) TemplateEngine {
	return TemplateEngine{templatesPath}
}

func (t *TemplateEngine) Parse(names ...string) *template.Template {
	files := []string{t.templatesPath + "/application.gohtml"}
	for _, name := range names {
		files = append(files, t.templatesPath+"/"+name+".gohtml")
	}
	return template.Must(template.ParseFiles(files...))
}

type ViewContext struct {
	w http.ResponseWriter
	r *http.Request
}

func (vc ViewContext) selectTemplateName() string {
	templateName := "application"
	if vc.r.Header.Get("Hx-Request") == "true" {
		templateName = "content"
	}
	return templateName
}

func (vc ViewContext) execute(t *template.Template, data any) error {
	vc.w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.ExecuteTemplate(vc.w, vc.selectTemplateName(), data)

	return err
}
