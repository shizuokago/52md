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
}

func getUser(r *http.Request) (User, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	key := datastore.NewKey(c, "User", u.ID, 0, nil)
	rtn := User{}
	if err := datastore.Get(c, key, &rtn); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		}
	}
	return rtn, nil
}

func putUser(r *http.Request) (User, error) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	r.ParseForm()
	rtn := User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		Url:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
	}
	_, err := datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &rtn)
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

func meRender(tName string, obj interface{}) {
	tmpl, err := template.ParseFiles("./templates/me/layout.tmpl", tName)
	if err != nil {
		return
	}
	err = tmpl.Execute(w, obj)
	if err != nil {
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	var u User
	if r.Method == "POST" {
		u, _ := putUser(r)
	} else {
		u, _ := getUser(r)
	}

	meRender("./templates/me/profile.tmpl", u)
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
	meRender("./templates/me/top.tmpl", nil)
}
