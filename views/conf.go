package views

import (
	"path/filepath"

	"github.com/elos/ehttp/templates"
)

const (
	Register templates.Name = iota
)

var (
	appPath               = "github.com/elos/mist"
	AssetsDir             = filepath.Join(templates.PackagePath(appPath), "assets")
	TemplatesDir          = filepath.Join(AssetsDir, "templates")
	layoutTemplate string = "layout.tmpl"
)

// Layout prepends variadic arguments with the layoutTemplate
func Layout(v ...string) []string {
	return templates.Prepend(layoutTemplate, v...)
}

var TemplatesSet = templates.TemplateSet{
	Register: Layout("register.tmpl"),
}
