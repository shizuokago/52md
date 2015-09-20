package go2md

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"net/http"
)

type User struct {
	UserKey   string
	Name      string
	Job       string
	Email     string
	Url       string
	TwitterId string
	LastWord  string
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

func getPublicUser(r *http.Request, userKey string) (*User, error) {

	c := appengine.NewContext(r)
	q := datastore.NewQuery("User").Filter("UserKey = ", userKey)
	var u []User
	q.GetAll(c, &u)

	return &u[0], nil
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
