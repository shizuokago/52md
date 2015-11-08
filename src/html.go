package gopredit

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strings"
	"time"
)

type Html struct {
	UserKey string
	Title   string
	Content []byte
	Date    time.Time
}

type HtmlJson struct {
	Success bool
	Html    *Html
}

func init() {
	http.HandleFunc("/me/slide/publish/", publishHandler)
}

// json access
func publishHandler(w http.ResponseWriter, r *http.Request) {

	//redfactoring
	editHandler(w, r)

	urls := strings.Split(r.URL.Path, "/")
	keyId := urls[4]

	u, err := getUser(r)
	if err != nil {
		//User not fount
	}

	s, err := getSlide(r, keyId)
	if err != nil {
		//Slide not fount
	}

	c := appengine.NewContext(r)

	id := u.UserKey + "/" + keyId
	key := createKey(c, "Html", id)

	data := Who{
		author:  u.UserKey,
		request: r,
	}

	content, err := createSlide(u, s, &data)
	if err != nil {
		//Slide Build Error
	}

	html := Html{
		UserKey: u.UserKey,
		Content: content,
		Title:   s.Title,
		Date:    time.Now(),
	}

	err = putHtml(c, key, &html)
	if err != nil {
		//Slide Publish Error
	}
}

func getHtml(c context.Context, key *datastore.Key) (*Html, error) {
	var h Html
	err := get(c, key, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func putHtml(c context.Context, key *datastore.Key, h *Html) error {
	_, err := datastore.Put(c, key, h)
	if err != nil {
		return err
	}
	return nil
}
