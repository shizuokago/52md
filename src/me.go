package go2md

import (
	"net/http"

	"html/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/me/", meHandler)
	http.HandleFunc("/me/profile", profileHandler)
	http.HandleFunc("/me/register", registerHandler)
}

func getUser(r *http.Request) (*User, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	key := datastore.NewKey(c, "User", u.ID, 0, nil)
	rtn := User{}
	if err := datastore.Get(c, key, &rtn); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return &rtn, nil
}

func putUser(r *http.Request) (*User, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	r.ParseForm()
	rtn := User{
		UserKey:   r.FormValue("UserKey"),
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		Url:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
		LastWord:  r.FormValue("LastWord"),
	}
	_, err := datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &rtn)
	if err != nil {
		return nil, err
	}
	return &rtn, nil
}

func meRender(w http.ResponseWriter, tName string, obj interface{}) {

	//select user slide

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
	c := appengine.NewContext(r)
	u := user.Current(c)
	r.ParseForm()
	rtn := User{
		UserKey: r.FormValue("UserKey"),
	}
	_, err := datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &rtn)
	if err != nil {
		panic(err)
	}
	meRender(w, "./templates/me/profile.tmpl", rtn)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	var u *User
	if r.Method == "POST" {
		u, _ = putUser(r)
	} else {
		u, _ = getUser(r)
	}
	meRender(w, "./templates/me/profile.tmpl", u)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	u := user.Current(c)

	if u == nil {
		url, _ := user.LoginURL(c, "/me/")
		http.Redirect(w, r, url, 301)
		return
	}

	du, err := getUser(r)
	if err != nil {
		panic(err)
	}

	if du == nil {
		meRender(w, "./templates/me/userkey.tmpl", nil)
	} else {

		//select user slide
		userkey := du.UserKey
		q := datastore.NewQuery("Slide").
			Filter("UserKey = ", userkey).
			Order("- SpeakDate")
		var s []Slide
		keys, _ := q.GetAll(c, &s)

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
