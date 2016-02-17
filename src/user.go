package gopredit

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"

	"net/http"
	"strconv"
)

type User struct {
	UserKey   string
	Name      string
	Job       string
	Email     string
	Url       string
	TwitterId string
	LastWord  string
	Size      int64
}

func existUser(r *http.Request, k string) bool {

	c := appengine.NewContext(r)
	if k == "me" || k == "file" {
		return true
	}

	q := datastore.NewQuery("User").Filter("UserKey = ", k)
	count, err := q.Count(c)
	if err != nil {
		return true
	}

	if count >= 1 {
		return true
	}
	return false
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

	size, err := strconv.ParseInt(r.FormValue("Size"), 10, 64)
	if err != nil {
		return nil, err
	}
	rtn := User{
		UserKey:   r.FormValue("UserKey"),
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		Url:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
		LastWord:  r.FormValue("LastWord"),
		Size:      size,
	}

	_, err = datastore.Put(c, datastore.NewKey(c, "User", u.ID, 0, nil), &rtn)
	if err != nil {
		return nil, err
	}
	return &rtn, nil
}
