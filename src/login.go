package go2md

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/user"
	"html/template"
)

func init() {
	http.HandleFunc("/me", meHandler)
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

	tmpl, err := template.ParseFiles("./templates/me.tmpl")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		panic(err)
	}

}
