package templates

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
)

//go:embed *.gohtml layouts/*.gohtml
var files embed.FS

type Templates struct {
	templates map[string]*template.Template
}

func Build() (*Templates, error) {
	templates := make(map[string]*template.Template)
	tmplFiles, err := fs.ReadDir(files, ".")
	if err != nil {
		return nil, fmt.Errorf("can't read templates: %w", err)
	}

	for _, tmplFile := range tmplFiles {
		log.Warnf("tmpl: %s", tmplFile.Name())

		if tmplFile.IsDir() {
			continue
		}

		tmpl, err := template.New(tmplFile.Name()).
			Funcs(sprig.FuncMap()).
			ParseFS(files, tmplFile.Name(), "layouts/*.gohtml")
		if err != nil {
			return nil, fmt.Errorf("can't parse template %s: %w", tmplFile.Name(), err)
		}

		templates[tmplFile.Name()] = tmpl.Funcs(sprig.FuncMap())
	}
	return &Templates{
		templates: templates,
	}, nil
}

func (ts *Templates) MustGet(name string) *template.Template {
	t, ok := ts.templates[name]
	if !ok {
		log.Fatalf("template %s not found", name)
	}
	return t
}
