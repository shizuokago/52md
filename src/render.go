package go2md

import (
	"golang.org/x/tools/present"
	"html/template"
	"path/filepath"
)

type Who struct {
	author string
	id     string
	data   string //no
}

func (s Who) AttributeFile(name string) ([]byte, error) {
	return []byte(s.data), nil
}

func createTemplate() (*template.Template, error) {
	base := "./"
	actionTmpl := filepath.Join(base, "templates/action.tmpl")
	contentTmpl := filepath.Join(base, "templates/slides.tmpl")
	tmpl := present.Template()
	tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
	if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}
