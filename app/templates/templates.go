package templates

import (
	"html/template"
	"io"
	"path/filepath"
	"filemanager/app/config"
)

type Template struct {
	Tmpl *template.Template
	Data      interface{}
	Writer    io.Writer
	Block	  interface{}
	// For import functionality
	Folder	interface{}
}

func Add(t ...string) *template.Template {
	return template.Must(template.ParseFiles(t...))
}

func (t *Template) Execute() error {
	if t.Block != nil {
		return t.Tmpl.ExecuteTemplate(t.Writer, t.Block.(string), t.Data)
	}
	return t.Tmpl.Execute(t.Writer, t.Data)
}

func (t *Template) AddWriter(w io.Writer) *Template {
	t.Writer = w
	return t
}

func (t *Template) Templates(templates []string) *Template {
	 t.Tmpl = template.Must(template.ParseFiles(templates...))
	 return t
}

func (t *Template) WithBlock(block interface{}) *Template {
	t.Block = block
	return t
}


// Importing the templates and adding path from the config.
func (t *Template) Import(templates ...string){
	var folder string = ""
	if t.Folder != nil {
		folder = t.Folder.(string)
	}
	for i, v := range templates {
		templates[i] = filepath.Join(config.TemplatesFolder, folder, v)
	}
	t.Tmpl = Add(templates...)
}

func (t *Template) FromFolder(folderName string) *Template {
	t.Folder = folderName
	return t
}