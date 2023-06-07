package http

import "html/template"

type TemplatesProvider interface {
    MustGet(name string) *template.Template
}
