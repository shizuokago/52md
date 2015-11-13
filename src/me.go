package gopredit

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
	"html/template"

	"net/http"
)

func init() {
	http.HandleFunc("/me/", meHandler)
	http.HandleFunc("/me/profile", profileHandler)
	http.HandleFunc("/me/register", registerHandler)
}

func meRender(w http.ResponseWriter, tName string, obj interface{}) {
	tmpl, err := template.ParseFiles("./templates/me/layout.tmpl", tName)
	if err != nil {
		return
	}
	err = tmpl.Execute(w, obj)
	if err != nil {
		return
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	// judge post

	c := appengine.NewContext(r)
	u := user.Current(c)
	r.ParseForm()

	userKey := r.FormValue("UserKey")
	if existUser(r, userKey) {
		return
	}

	rtn := User{
		UserKey: userKey,
	}

	_, err := datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &rtn)
	if err != nil {
	}
	//Profile Page
	meRender(w, "./templates/me/profile.tmpl", rtn)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	var u *User
	var err error
	// add error handling
	if r.Method == "POST" {
		u, err = putUser(r)
	} else {
		u, err = getUser(r)
	}

	if err != nil {
	}

	meRender(w, "./templates/me/profile.tmpl", u)
}

func meHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		url, _ := user.LoginURL(c, "/me/")
		http.Redirect(w, r, url, 301)
		return
	}

	du, err := getUser(r)
	if err != nil {
	}

	if du == nil {
		//register userkey
		meRender(w, "./templates/me/userkey.tmpl", nil)
	} else {
		//select user slide
		userkey := du.UserKey
		q := datastore.NewQuery("Slide").
			Filter("UserKey = ", userkey).
			Order("- SpeakDate")
		var s []Slide

		// function get user Slide

		keys, err := q.GetAll(c, &s)
		if err != nil {
		}

		rtn := make([]TemplateSlide, len(s))
		for i, elm := range s {
			rtn[i] = TemplateSlide{
				Title:     elm.Title,
				SubTitle:  elm.SubTitle,
				SpeakDate: elm.SpeakDate,
				Key:       keys[i].StringID(),
			}
		}
		meRender(w, "./templates/me/top.tmpl", rtn)
	}
}

type TemplateSlide struct {
	Title     string
	SubTitle  string
	SpeakDate string
	Key       string
}
