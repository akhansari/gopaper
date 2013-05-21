// post table

package models

import (
	"appengine"
	"appengine/datastore"
	"html"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gopaper/ex/stringsex"
)

type Post struct {
	Id       int64 `datastore:"-"`
	Title    string
	Body     string `datastore:",noindex"`
	Publish  bool
	Page     bool
	AddDate  time.Time
	EditDate time.Time
	Tags     string
}

// number of posts per page
const (
	postPerPage = 3
)

// save function for PropertyLoadSaver interface
func (p *Post) save(c chan<- datastore.Property) error {
	defer close(c)
	return datastore.SaveStruct(p, c)
}

// add the post
func (p *Post) Add(context appengine.Context) error {
	return datastore.RunInTransaction(context, func(context appengine.Context) error {

		p.AddDate = time.Now()

		postKey, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Post", nil), p)
		if err != nil {
			return err
		}

		if p.Publish {
			err = addTags(context, postKey, p.Tags, p.AddDate)
			if err != nil {
				return err
			}
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})
}

// edit the post
func (p *Post) Edit(context appengine.Context) error {
	return datastore.RunInTransaction(context, func(context appengine.Context) error {

		p.EditDate = time.Now()

		postKey := datastore.NewKey(context, "Post", "", p.Id, nil)
		_, err := datastore.Put(context, postKey, p)
		if err != nil {
			return err
		}

		deleteTags(context, postKey)
		if p.Publish {
			err = addTags(context, postKey, p.Tags, p.AddDate)
			if err != nil {
				return err
			}
		}

		return nil
	}, &datastore.TransactionOptions{XG: true})
}

// fill the post from html fields
func (p *Post) FeedFromForm(form url.Values) {
	p.Title = form.Get("Title")
	p.Body = html.UnescapeString(form.Get("Body"))
	p.Publish = len(form.Get("Publish")) > 0
	p.Page = len(form.Get("Page")) > 0
	p.Tags = strings.Join(stringsex.SplitAndTrimSpace(form.Get("Tags"), ","), ",")
}

// get the post from its id
func GetPostFromStrId(context appengine.Context, strId string) (*Post, error) {

	var post Post
	id, _ := strconv.ParseInt(strId, 10, 64)
	postKey := datastore.NewKey(context, "Post", "", id, nil)
	err := datastore.Get(context, postKey, &post)
	post.Id = id
	return &post, err
}

// delete the post from its id
func DeletePostFromStrId(context appengine.Context, strId string) error {

	id, _ := strconv.ParseInt(strId, 10, 64)
	postKey := datastore.NewKey(context, "Post", "", id, nil)
	err := deleteTags(context, postKey)
	if err != nil {
		return err
	}
	return datastore.Delete(context, postKey)
}

// get all posts from a query
func GetPostsFromQuery(context appengine.Context, query *datastore.Query) (*[]Post, error) {

	var posts []Post
	keys, err := query.GetAll(context, &posts)
	if err != nil {
		return nil, err
	}
	for i, key := range keys {
		posts[i].Id = key.IntID()
	}
	return &posts, nil
}

// get posts from a query per page
func GetPostsFromQueryPerPage(context appengine.Context, query *datastore.Query, page int) (*[]Post, bool, error) {

	cnt, _ := query.Count(context)
	offset := (page - 1) * postPerPage
	hasOlder := offset+postPerPage < cnt

	qry := query.Limit(postPerPage).Offset(offset)

	posts, err := GetPostsFromQuery(context, qry)

	return posts, hasOlder, err
}

// get all published posts
func GetPublishedPostsPerPage(context appengine.Context, page int) (*[]Post, bool, error) {

	qry := datastore.NewQuery("Post").
		Filter("Publish =", true).Filter("Page =", false).
		Order("-AddDate")

	posts, hasOlder, err := GetPostsFromQueryPerPage(context, qry, page)

	return posts, hasOlder, err
}

// get all published posts of a tag
func GetTagsPerPage(context appengine.Context, tag string, page int) (*[]Post, bool, error) {

	qry := datastore.NewQuery("Tag").KeysOnly().
		Filter("Name =", tag).Order("-PostDate")

	cnt, _ := qry.Count(context)
	offset := (page - 1) * postPerPage
	hasOlder := offset+postPerPage < cnt

	qry = qry.Limit(postPerPage).Offset(offset)

	tagKeys, err := qry.GetAll(context, nil)

	posts := make([]Post, len(tagKeys))
	for i := 0; i < len(tagKeys); i++ {
		postKey := tagKeys[i].Parent()
		datastore.Get(context, postKey, &posts[i])
		posts[i].Id = postKey.IntID()
	}

	return &posts, hasOlder, err
}
