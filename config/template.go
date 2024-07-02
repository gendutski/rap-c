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
	layoutPath   string = "layout.html"
)

type Renderer struct {
	autoReloadTemplate bool // not for production
	keys               map[string][]string
	templates          map[string]*Templates
}

type Templates struct {
	Template        *template.Template
	ExecuteTemplate string
}

func NewRenderer(autoReloadTemplate bool) (*Renderer, error) {
	keys := map[string][]string{
		"404.html":   {"404.html"},
		"500.html":   {"500.html"},
		"login.html": {"login.html"},
	}

	templates, err := loadTemplates(keys)
	if err != nil {
		return nil, err
	}
	return &Renderer{autoReloadTemplate: autoReloadTemplate, keys: keys, templates: templates}, nil
}

func loadTemplates(keys map[string][]string) (map[string]*Templates, error) {
	templates := make(map[string]*Templates)
	for key, values := range keys {
		if len(values) == 0 {
			continue
		}

		var files []string
		for _, v := range values {
			files = append(files, filepath.Join(storagePath, templatePath, v))
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		templates[key] = &Templates{
			Template:        tmpl,
			ExecuteTemplate: values[0],
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
		err = tmpl.Template.ExecuteTemplate(w, tmpl.ExecuteTemplate, data)
		if err != nil {
			return fmt.Errorf("failed render template: %s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("template %s not found", name)
}
