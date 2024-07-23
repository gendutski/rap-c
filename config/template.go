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
	keys               map[string]*TemplateMaper
	templates          map[string]*Templates
}

type Templates struct {
	Template        *template.Template
	ExecuteTemplate string
}

type TemplateMaper struct {
	Files           []string
	ExecuteTemplate string
}

var funcMap = template.FuncMap{
	"partialExists": func(name string, t *template.Template) bool {
		return t.Lookup(name) != nil
	},
}

func NewRenderer(autoReloadTemplate bool) (*Renderer, error) {
	keys := map[string]*TemplateMaper{
		"401.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "401.html"),
			},
			ExecuteTemplate: "index",
		},
		"404.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "404.html"),
			},
			ExecuteTemplate: "index",
		},
		"error.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "error.html"),
			},
			ExecuteTemplate: "index",
		},
		"login.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "login.html"),
			},
			ExecuteTemplate: "index",
		},
		"pass-changer.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "pass-changer.html"),
			},
			ExecuteTemplate: "index",
		},
		"forgot-password.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "forgot-password.html"),
			},
			ExecuteTemplate: "index",
		},
		"reset-password.html": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "reset-password.html"),
			},
			ExecuteTemplate: "index",
		},
		"profile": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "layouts", "layout.html"),
				filepath.Join(storagePath, templatePath, "layouts", "sidebar-menu.html"),
				filepath.Join(storagePath, templatePath, "profile", "index.html"),
			},
			ExecuteTemplate: "layout",
		},
		"dashboard": {
			Files: []string{
				filepath.Join(storagePath, templatePath, "layouts", "layout.html"),
				filepath.Join(storagePath, templatePath, "layouts", "sidebar-menu.html"),
				filepath.Join(storagePath, templatePath, "dashboard", "index.html"),
			},
			ExecuteTemplate: "layout",
		},
	}

	templates, err := loadTemplates(keys)
	if err != nil {
		return nil, err
	}
	return &Renderer{autoReloadTemplate: autoReloadTemplate, keys: keys, templates: templates}, nil
}

func loadTemplates(keys map[string]*TemplateMaper) (map[string]*Templates, error) {
	templates := make(map[string]*Templates)
	for key, values := range keys {
		tmpl, err := template.New(values.ExecuteTemplate).
			Funcs(funcMap).
			ParseFiles(values.Files...)
		if err != nil {
			return nil, err
		}

		templates[key] = &Templates{
			Template:        tmpl,
			ExecuteTemplate: values.ExecuteTemplate,
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
		// err = tmpl.Template.Execute(w, data)
		if err != nil {
			return fmt.Errorf("failed render template: %s", err.Error())
		}
		return nil
	}
	return fmt.Errorf("template %s not found", name)
}
