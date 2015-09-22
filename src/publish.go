package go2md

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
	Content []byte
	Date    time.Time
}

type HtmlJson struct {
	Success bool
	Html    *Html
}

func init() {
	http.HandleFunc("/", publicHandler)
	http.HandleFunc("/me/slide/publish/", publishHandler)
}

// json access
func publishHandler(w http.ResponseWriter, r *http.Request) {

	urls := strings.Split(r.URL.Path, "/")
	keyId := urls[4]

	u, err := getUser(r)
	if err != nil {
	}

	s, err := getSlide(r, keyId)
	if err != nil {
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
	}

	html := Html{
		UserKey: u.UserKey,
		Content: content,
		Date:    time.Now(),
	}

	err = putHtml(c, key, &html)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {

	urls := strings.Split(r.URL.Path, "/")
	userKey := urls[1]
	keyId := urls[2]
	id := userKey + "/" + keyId

	c := appengine.NewContext(r)

	key := createKey(c, "Html", id)
	html, err := getHtml(c, key)
	if err != nil {
		keyName := userKey + "/" + strings.Join(urls[2:], "/")
		f, _ := getFile(r, keyName)
		if f != nil {
			w.Write(f.Data)
		} else {
		}
		return

	}

	_, err = w.Write(html.Content)
	if err != nil {
		//err page
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
