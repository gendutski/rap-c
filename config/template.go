package config

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/labstack/echo/v4"
)

const (
	storagePath  string = "storage"
	templatePath string = "templates"
)

type Renderer struct {
	autoReloadTemplate bool // not for production
	keys               map[string][][]string
	templates          map[string]*Templates
}

type Templates struct {
	Template        *template.Template
	ExecuteTemplate string
}

func NewRenderer(autoReloadTemplate bool) (*Renderer, error) {
	keys := map[string][][]string{
		"401.html":            {{"401.html"}},
		"404.html":            {{"404.html"}},
		"error.html":          {{"error.html"}},
		"login.html":          {{"login.html"}},
		"pass-changer.html":   {{"pass-changer.html"}},
		"request-reset.html":  {{"request-reset.html"}},
		"reset-password.html": {{"reset-password.html"}},
		"profile.html":        {{"layouts", "layout.html"}, {"layouts", "sidebar-menu.html"}, {"profile", "index.html"}},
		"dashboard.html":      {{"layouts", "layout.html"}, {"layouts", "sidebar-menu.html"}, {"dashboard", "index.html"}},
	}

	templates, err := loadTemplates(keys)
	if err != nil {
		return nil, err
	}
	return &Renderer{autoReloadTemplate: autoReloadTemplate, keys: keys, templates: templates}, nil
}

func loadTemplates(keys map[string][][]string) (map[string]*Templates, error) {
	templates := make(map[string]*Templates)
	for key, values := range keys {
		if len(values) == 0 {
			continue
		}

		var files []string
		for _, v := range values {
			_path := []string{storagePath, templatePath}
			_path = append(_path, v...)
			files = append(files, filepath.Join(_path...))
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		templates[key] = &Templates{
			Template:        tmpl,
			ExecuteTemplate: values[0][len(values[0])-1:][0], // first template will get execute
		}
	}
	return templates, nil
}

func (r *Renderer) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	var err error

	if r.autoReloadTemplate {
		r.templates, err = loadTemplates(r.keys)
		if err != nil {
			return fmt.Errorf("failed when reload templates: %s", err.Error())
		}
	}

	if tmpl, ok := r.templates[name]; ok {
		// err = tmpl.Template.ExecuteTemplate(w, tmpl.ExecuteTemplate, data)
		err = tmpl.Template.Execute(w, data)
		if err != nil {
			return fmt.Errorf("failed render template: %s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("template %s not found", name)
}
