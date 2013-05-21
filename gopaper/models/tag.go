package models

import (
	"appengine"
	"appengine/datastore"
	"strings"
	"time"
)

type Tag struct {
	Name     string
	PostDate time.Time
}

func addTags(context appengine.Context, postKey *datastore.Key, tags string, postDate time.Time) error {

	tagsArr := strings.Split(tags, ",")
	for _, value := range tagsArr {
		if len(value) > 0 {
			tag := Tag{Name: value, PostDate: postDate}
			_, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Tag", postKey), &tag)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteTags(context appengine.Context, postKey *datastore.Key) error {
	tagsKeys, _ := datastore.NewQuery("Tag").Ancestor(postKey).KeysOnly().GetAll(context, nil)
	return datastore.DeleteMulti(context, tagsKeys)
}
