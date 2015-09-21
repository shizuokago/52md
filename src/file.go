package go2md

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"io/ioutil"
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/me/file/upload", uploadHandler)
}

type File struct {
	UserKey string
	Data    []byte
}

//change ajax access
func uploadHandler(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("fileName")
	file, _, err := r.FormFile("uploadFile")
	if err != nil {
		//add error handling
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
	http.Redirect(w, r, r.FormValue("redirect"), 301)
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

func getFileKey(r *http.Request) ([]string, error) {
	c := appengine.NewContext(r)
	u, _ := getUser(r)
	userKey := u.UserKey

	q := datastore.NewQuery("File").KeysOnly().Filter("UserKey = ", userKey)
	keys, _ := q.GetAll(c, nil)

	rtn := make([]string, len(keys))

	for idx, elm := range keys {
		rtn[idx] = strings.Replace(elm.StringID(), userKey, "", 1)
	}
	return rtn, nil

}
