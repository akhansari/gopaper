// backend controller

package handlers

import (
	"appengine/datastore"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"

	"gopaper/models"
)

type Backend struct {
	*Controller
}

// home
func (c Backend) Index() {

	vars := mux.Vars(c.Request)
	page, _ := strconv.Atoi(vars["id"])
	if page < 1 {
		page = 1
	}

	query := datastore.NewQuery("Post").Order("-AddDate")
	posts, hasOlder, _ := models.GetPostsFromQueryPerPage(c.Context, query, page)

	data := struct {
		Posts    *[]models.Post
		HasOlder bool
		HasNewer bool
		Page     int
		BaseUrl  string
	}{
		posts,
		hasOlder,
		page > 1 && posts != nil,
		page,
		"/backend",
	}

	c.NestedViews = []string{"Shared/PaginationNav"}
	c.Render(&data)
}

// adding new post
func (c Backend) AddPost() {

	post := new(models.Post)
	var errMsg string

	c.Request.ParseForm()
	if len(c.Request.Form) > 0 {

		post.FeedFromForm(c.Request.Form)
		err := post.Add(c.Context)
		if err == nil {
			c.Redirect("/backend")
			return
		} else {
			errMsg = err.Error()
		}
	}

	data := struct {
		Post   *models.Post
		ErrMsg string
	}{
		post,
		errMsg,
	}

	c.ViewName = "Backend/Post"
	c.Render(&data)
}

// showing and editing an existing post
func (c Backend) Post() {

	vars := mux.Vars(c.Request)
	var errMsg string

	post, err := models.GetPostFromStrId(c.Context, vars["id"])
	if err != nil {
		http.Error(c.Response, "Post not found", http.StatusNotFound)
		return
	}

	c.Request.ParseForm()
	if len(c.Request.Form) > 0 {

		post.FeedFromForm(c.Request.Form)
		err := post.Edit(c.Context)
		if err != nil {
			errMsg = err.Error()
		}
	}

	data := struct {
		Post   *models.Post
		ErrMsg string
	}{
		post,
		errMsg,
	}

	c.Render(&data)
}

// deleting a post
func (c Backend) DeletePost() {

	vars := mux.Vars(c.Request)
	err := models.DeletePostFromStrId(c.Context, vars["id"])
	if err != nil {
		http.Error(c.Response, err.Error(), http.StatusInternalServerError)
	} else {
		c.Redirect("/backend")
	}
}
