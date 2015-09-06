package go2md

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/godoc/static"
	"golang.org/x/tools/present"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

type Who struct {
	author string
	id     string
	data   string //no
}

func (s Who) AttributeFile(name string) ([]byte, error) {
	return []byte(s.data), nil
}

var scripts = []string{"jquery.js", "jquery-ui.js", "playground.js", "play.js"}
var modTime = time.Now()
var scriptByte []byte

func init() {
	playScript("./", "HTTPTransport")
	present.PlayEnabled = true
	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")
	http.HandleFunc("/play.js", playHandler)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/javascript")
	http.ServeContent(w, r, "", modTime, bytes.NewReader(scriptByte))
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

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
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
