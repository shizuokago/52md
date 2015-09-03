package go2md

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
	"strings"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"github.com/pborman/uuid"
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
	http.HandleFunc("/play.js", playHandler)

	http.HandleFunc("/me/slide/create", createHandler)
	http.HandleFunc("/me/slide/edit/", editHandler)
	http.HandleFunc("/me/slide/view/", viewHandler)
}

func createHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	// get user data
	u, _ := getUser(r)
	slide := Slide{
		UserKey:   u.UserKey,
		Title:     "EmptyTitle",
		SubTitle:  "EmptySubTitle",
		SpeakDate: "1 Sep 2015",
		Tags:      "golang,present",
		Markdown:  "* Page 1",
	}

	// add empty slide data
	key, _ := datastore.Put(c, datastore.NewKey(c, "Slide", uuid.New(), 0, nil), &slide)
	http.Redirect(w, r, "/me/slide/edit/"+key.StringID(), 301)
}

func putSlide(r *http.Request, key string) (*Slide, error) {
	c := appengine.NewContext(r)
	r.ParseForm()

	slide := Slide{
		UserKey:   r.FormValue("UserKey"),
		Title:     r.FormValue("Title"),
		SubTitle:  r.FormValue("SubTitle"),
		SpeakDate: r.FormValue("SpeakDate"),
		Tags:      r.FormValue("Tags"),
		Markdown:  r.FormValue("Markdown"),
	}
	datastore.Put(c, datastore.NewKey(c, "Slide", key, 0, nil), &slide)
	return &slide, nil
}

func getSlide(r *http.Request, key string) (*Slide, error) {
	c := appengine.NewContext(r)
	k := datastore.NewKey(c, "Slide", key, 0, nil)
	rtn := Slide{}
	if err := datastore.Get(c, k, &rtn); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return &rtn, nil
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	//c := appengine.NewContext(r)
	urls := strings.Split(r.URL.Path, "/")
	keyId := urls[4]

	var s *Slide
	if r.Method == "POST" {
		s, _ = putSlide(r, keyId)
	} else {
		s, _ = getSlide(r, keyId)
	}
	rtn := SlideView{
		Key:  keyId,
		Data: s,
	}
	meRender(w, "./templates/me/edit.tmpl", rtn)
}

type SlideView struct {
	Key  string
	Data *Slide
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	u, _ := getUser(r)

	urls := strings.Split(r.URL.Path, "/")
	keyId := urls[4]
	s, _ := getSlide(r, keyId)

	data := Who{
		author: "secondarykey",
		id:     "1",
	}

	slideTxt := ""
	slideTxt += s.Title + "\n"
	slideTxt += s.SubTitle + "\n"
	slideTxt += s.SpeakDate + "\n"
	slideTxt += "Tags:" + s.Tags + "\n"
	slideTxt += "\n"
	slideTxt += u.Name + "\n"
	slideTxt += u.Job + "\n"
	slideTxt += u.Url + "\n"
	slideTxt += "@" + u.TwitterId + "\n"
	slideTxt += "\n"
	slideTxt += s.Markdown

	c := appengine.NewContext(r)
	log.Infof(c, slideTxt)

	//52md
	//Golang Present Tools Editor
	//15 Aug 2015
	//Tags: golang shizuoka_go
	//
	//secondarykey
	//Programer
	//http://github.com/shizuokago/52md
	//@secondarykey
	//
	//* This Service Alpha

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
