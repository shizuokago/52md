package go2md

import (
	"golang.org/x/tools/present"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func init() {
	http.HandleFunc("/example", exampleHandler)
	http.HandleFunc("/example/change", changeExampleHandler)
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/example.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

type Who struct {
	author string
	id     string
	data   string //no
}

func (s Who) AttributeFile(name string) ([]byte, error) {
	return []byte(s.data), nil
}

func changeExampleHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	slideTxt := r.FormValue("slide")
	data := Who{
		author: "secondarykey",
		id:     "1",
	}

	ctx := present.Context{ReadFile: data.AttributeFile}
	reader := strings.NewReader(slideTxt)
	doc, err := ctx.Parse(reader, "tour.slide", 0)
	if err != nil {
		panic(err)
	}

	tmpl, err := createTemplate()
	if err != nil {
		panic(err)
	}
	doc.Render(w, tmpl)
}

func createTemplate() (*template.Template, error) {
	base := "./"
	actionTmpl := filepath.Join(base, "templates/action.tmpl")
	contentTmpl := filepath.Join(base, "templates/slides_body.tmpl")
	tmpl := present.Template()
	tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
	if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}
