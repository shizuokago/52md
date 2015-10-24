package gopredit

import (
	"google.golang.org/appengine"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/", publicHandler)
}

func publicHandler(w http.ResponseWriter, r *http.Request) {

	urls := strings.Split(r.URL.Path, "/")
	userKey := urls[1]
	keyId := urls[2]

	if keyId == "file" {
		keyName := userKey + "/" + strings.Join(urls[3:], "/")
		f, _ := getFile(r, keyName)
		if f != nil {
			w.Write(f.Data)
		} else {
		}
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
		//err page

	}
}
