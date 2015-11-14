package gopredit

import (
	"net/http"
)

type Who struct {
	author  string
	request *http.Request
}

func (s Who) AttributeFile(name string) ([]byte, error) {
	key := s.author + "/" + name
	f, err := getFile(s.request, key)
	if err != nil {
		return nil,err
	}
	if f == nil {
		return []byte("Not Found"), nil
	}
	return f.Data, nil
}
