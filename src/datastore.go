package gopredit

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func init() {
}

func createKey(c context.Context, kind, id string) *datastore.Key {
	key := datastore.NewKey(c, kind, id, 0, nil)
	return key
}

func get(c context.Context, key *datastore.Key, v interface{}) error {
	if err := datastore.Get(c, key, v); err != nil {
		return err
	}
	return nil
}

func put(c context.Context, key *datastore.Key, v interface{}) error {
	_, err := datastore.Put(c, key, v)
	if err != nil {
		return err
	}

	return nil
}
