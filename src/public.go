package gopredit

import (
	"golang.org/x/tools/present"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"html/template"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/", publicHandler)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {

	urls := strings.Split(r.URL.Path, "/")
	count := len(urls)
	if count < 3 {
		errorPage(w, "URL Error", "Argument Error", "", 400)
		return
	}

	userKey := urls[1]
	keyId := urls[2]
	if keyId == "" {
		// User Page
		renderUserPage(w, r, userKey)
		return
	}

	if keyId == "file" {
		keyName := userKey + "/" + strings.Join(urls[3:], "/")

		f, _ := getFile(r, keyName)
		if f != nil {
			w.Write(f.Data)
		} else {
			errorPage(w, "Error", "Not Found", "This is not found", 404)
		}
		return
	}

	if count > 3 {
		errorPage(w, "URL Error", "Argument Error", "", 400)
		return
	}

	id := userKey + "/" + keyId
	c := appengine.NewContext(r)

	key := createKey(c, "Html", id)
	html, err := getHtml(c, key)
	if err != nil {
		errorPage(w, "Error", "Not Found", "This is not found", 404)
		return
	}

	//response write
	_, err = w.Write(html.Content)
	if err != nil {
		errorPage(w, "Error", "HTML Write Error", "", 500)
	}
}

func renderUserPage(w http.ResponseWriter, r *http.Request, userKey string) {

	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("UserKey = ", userKey)
	var us []User
	_, err := q.GetAll(c, &us)
	if err != nil {
		errorPage(w, "InternalServerError", "User Error", "", 500)
		return
	}
	if len(us) == 0 {
		errorPage(w, "NotFound", "User NotFound", "", 404)
		return
	}
	var hs []Html
	qh := datastore.NewQuery("Html").Filter("UserKey = ", userKey)
	keys, err := qh.GetAll(c, &hs)
	if err != nil {
		errorPage(w, "InternalServerError", "Html Error", "", 500)
		return
	}
	if len(hs) == 0 {
		errorPage(w, "NotFound", "User NotFound", "", 404)
		return
	}

	u := us[0]

	slideTxt := ""
	slideTxt = addLine(slideTxt, u.Name+" Slides", "")
	slideTxt = addLine(slideTxt, "golang,GoPreEdit", "Tags:")
	slideTxt += "\n"

	slideTxt = addLine(slideTxt, u.Name, "")
	slideTxt = addLine(slideTxt, u.Job, "")
	slideTxt = addLine(slideTxt, u.Url, "")
	slideTxt = addLine(slideTxt, u.TwitterId, "@")

	slideTxt += "\n"

	for idx, elm := range hs {
		slideTxt = addLine(slideTxt, "* "+elm.Title, "")
		slideTxt = addLine(slideTxt, ".link ../"+keys[idx].StringID()+" "+elm.Title, "")
		slideTxt += "\n"
	}
	data := Who{
		author:  userKey,
		request: r,
	}

	ctx := present.Context{ReadFile: data.AttributeFile}
	reader := strings.NewReader(slideTxt)
	doc, err := ctx.Parse(reader, "tour.slide", 0)
	if err != nil {
		errorPage(w, "InternalServerError", "Parse", "", 500)
		return
	}

	tmpl, err := createTemplate()
	if err != nil {
		errorPage(w, "InternalServerError", "create template error", "", 500)
		return
	}

	if u.LastWord == "" {
		u.LastWord = "Thank you"
	}

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		LastWord    string
	}{doc, tmpl, true, u.LastWord}

	err = tmpl.ExecuteTemplate(w, "root", rtn)
	if err != nil {
		errorPage(w, "InternalServerError", "execute template error", "", 500)
	}
}
