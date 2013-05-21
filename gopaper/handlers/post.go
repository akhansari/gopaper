// post controller

package handlers

import (
	"github.com/gorilla/mux"
	"net/http"

	"gopaper/models"
)

type Post struct {
	*Controller
}

func (c Post) Index() {

	vars := mux.Vars(c.Request)

	post, err := models.GetPostFromStrId(c.Context, vars["id"])
	if err != nil || !post.Publish {
		http.Error(c.Response, "Post not found", http.StatusNotFound)
		return
	}
	postShow := models.FormatPost(post)

	data := struct {
		Post *models.PostShow
	}{
		postShow,
	}

	c.NestedViews = []string{"Shared/PostTagsList"}
	c.Render(&data)

}
