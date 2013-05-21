// homepage controller

package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"gopaper/models"
)

type Home struct {
	*Controller
}

func (c Home) Index() {

	vars := mux.Vars(c.Request)
	page, _ := strconv.Atoi(vars["id"])
	if page < 1 {
		page = 1
	}

	posts, hasOlder, _ := models.GetPublishedPostsPerPage(c.Context, page)
	postsShow := models.FormatPosts(posts)

	data := struct {
		Posts    *[]models.PostShow
		HasOlder bool
		HasNewer bool
		Page     int
		BaseUrl  string
	}{
		postsShow,
		hasOlder,
		page > 1 && posts != nil,
		page,
		"",
	}

	c.NestedViews = []string{"Shared/PostTagsList", "Shared/PaginationNav"}
	c.Render(&data)

}

func (c Home) Tags() {

	vars := mux.Vars(c.Request)

	tag, _ := vars["tag"]

	page, _ := strconv.Atoi(vars["id"])
	if page < 1 {
		page = 1
	}

	posts, hasOlder, _ := models.GetTagsPerPage(c.Context, tag, page)
	postsShow := models.FormatPosts(posts)

	data := struct {
		Posts    *[]models.PostShow
		HasOlder bool
		HasNewer bool
		Page     int
		BaseUrl  string
	}{
		postsShow,
		hasOlder,
		page > 1 && posts != nil,
		page,
		"",
	}

	c.ViewName = "Home/Index"
	c.NestedViews = []string{"Shared/PostTagsList", "Shared/PaginationNav"}
	c.Render(&data)

}
