package go2md

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"net/http"
)

func init() {
	http.HandleFunc("/me/file/upload", uploadHandler)
}

type File struct {
	UserKey string
	Data    []byte
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("fileName")
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)

	c := appengine.NewContext(r)
	// get user data
	u, _ := getUser(r)
	key := datastore.NewKey(c, "File", u.UserKey+"/"+name, 0, nil)
	f := File{
		UserKey: u.UserKey,
		Data:    b,
	}
	// add empty slide data
	datastore.Put(c, key, &f)
}

func getFile(r *http.Request, name string) (*File, error) {
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "File", name, 0, nil)
	rtn := File{}
	if err := datastore.Get(c, key, &rtn); err != nil {
		if err != datastore.ErrNoSuchEntity {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return &rtn, nil
}
