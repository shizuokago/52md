package shizuokago

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/godoc/static"
	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
	"html/template"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	contentTemplate map[string]*template.Template
)
var scripts = []string{"jquery.js", "jquery-ui.js", "playground.js", "play.js"}

func init() {
	basePath := "./"
	initTemplates(basePath)
	playScript(basePath, "HTTPTransport")
	present.PlayEnabled = true
	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")
	http.HandleFunc("/slides/", slideHandler)
	http.HandleFunc("/play.js", playHandler)
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func slideHandler(w http.ResponseWriter, r *http.Request) {
	const base = "./slides"
	name := filepath.Join(base, r.URL.Path)
	if isDoc(name) {
		err := renderDoc(w, name)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	http.FileServer(http.Dir(base)).ServeHTTP(w, r)
}

func isDoc(path string) bool {
	_, ok := contentTemplate[filepath.Ext(path)]
	return ok
}

func initTemplates(base string) error {
	// Locate the template file.
	actionTmpl := filepath.Join(base, "templates/action.tmpl")

	contentTemplate = make(map[string]*template.Template)

	for ext, contentTmpl := range map[string]string{
		".slide": "slides.tmpl",
	} {
		contentTmpl = filepath.Join(base, "templates", contentTmpl)

		// Read and parse the input.
		tmpl := present.Template()
		tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
		if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
			return err
		}
		contentTemplate[ext] = tmpl
	}
	return nil
}

// renderDoc reads the present file, gets its template representation,
// and executes the template, sending output to w.
func renderDoc(w io.Writer, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := parse(docFile, 0)
	if err != nil {
		return err
	}
	// Find which template should be executed.
	tmpl := contentTemplate[filepath.Ext(docFile)]
	// Execute the template.
	return doc.Render(w, tmpl)
}

func parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return present.Parse(f, name, 0)
}

var modTime = time.Now()
var scriptByte []byte

func playHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/javascript")
	http.ServeContent(w, r, "", modTime, bytes.NewReader(scriptByte))
}

func playScript(root, transport string) {
	var buf bytes.Buffer
	for _, p := range scripts {
		if s, ok := static.Files[p]; ok {
			buf.WriteString(s)
			continue
		}
		b, err := ioutil.ReadFile(filepath.Join(root, "./static", p))
		if err != nil {
			panic(err)
		}
		buf.Write(b)
	}
	fmt.Fprintf(&buf, "\ninitPlayground(new %v());\n", transport)
	scriptByte = buf.Bytes()
}
