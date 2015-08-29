package go2md

import (
	"fmt"
	"net/http"

	"html/template"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/me", meHandler)
	http.HandleFunc("/me/profile", profileHandler)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	gu := user.Current(c)

	if r.Method == "POST" {
		r.ParseForm()
		u := User{
			Name:      r.FormValue("Name"),
			Job:       r.FormValue("Job"),
			Email:     r.FormValue("Email"),
			Url:       r.FormValue("Url"),
			TwitterId: r.FormValue("TwitterId"),
		}
		datastore.Put(c, datastore.NewKey(c, "User", gu.ID, 0, nil), &u)

		//err = datastore.RunInTransaction(c, func(c context.Context) error {
		//_, err = datastore.Put(c, key, &u)
		//})
	}

	key := datastore.NewKey(c, "User", gu.ID, 0, nil)
	u := User{}
	if err := datastore.Get(c, key, &u); err != nil {
		if err != datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	tmpl, err := template.ParseFiles("./templates/me/me_layout.tmpl", "./templates/me/profile.tmpl")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, u)
	if err != nil {
		panic(err)
	}
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ := user.LoginURL(c, "/me")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	//url, _ :=user.LogoutURL(c, "/logout")

	tmpl, err := template.ParseFiles("./templates/me/me_layout.tmpl", "./templates/me/me.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}

}
