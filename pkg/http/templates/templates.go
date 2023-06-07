package templates

import (
    "embed"
    "fmt"
    "github.com/Masterminds/sprig/v3"
    log "github.com/sirupsen/logrus"
    "html/template"
    "io/fs"
)

var (
    //go:embed *.gohtml layouts/*.gohtml
    files embed.FS
)

type Templates struct {
    templates map[string]*template.Template
}

type Provider interface {
    MustGet(name string) *template.Template
}

func Build() (*Templates, error) {
    ts := make(map[string]*template.Template)
    tmplFiles, err := fs.ReadDir(files, ".")
    if err != nil {
        return nil, err
    }

    for _, tmpl := range tmplFiles {
        log.Warnf("tmpl: %s", tmpl.Name())

        if tmpl.IsDir() {
            continue
        }

        pt, err := template.New(tmpl.Name()).
            Funcs(sprig.FuncMap()).
            ParseFS(files, tmpl.Name(), "layouts/*.gohtml")
        if err != nil {
            return nil, fmt.Errorf("can't parse template %s: %v", tmpl.Name(), err)
        }

        ts[tmpl.Name()] = pt.Funcs(sprig.FuncMap())
    }
    return &Templates{
        templates: ts,
    }, nil
}

func (ts *Templates) MustGet(name string) *template.Template {
    t, ok := ts.templates[name]
    if !ok {
        log.Fatalf("template %s not found", name)
    }
    return t
}
