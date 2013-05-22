package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopaper/ex/stringsex"
)

type PostShow struct {
	*Post
	Day       string
	Month     string
	TagsArr   []string
	Url       string
	ShortDesc string
}

// prepare the stuff for display
func FormatPosts(posts *[]Post) *[]PostShow {

	postsShow := make([]PostShow, len(*posts))

	for i, post := range *posts {

		// take only the first paragraph
		rgx, _ := regexp.Compile("<p>(.*?)</p>")
		if rgx.MatchString(post.Body) {
			post.Body = rgx.FindString(post.Body)
		}

		postsShow[i] = *FormatPost(post)
	}
	return &postsShow
}

// prepare the post for display
func FormatPost(post Post) *PostShow {

	postShow := PostShow{
		&post,
		fmt.Sprintf("%02d", post.AddDate.Day()),
		post.AddDate.Month().String()[:3],
		strings.Split(post.Tags, ","),
		"/post/" + strconv.FormatInt(post.Id, 10) + "/" + stringsex.FormatUrl(post.Title),
		"",
	}

	if len(post.Body) > 150 {
		postShow.ShortDesc = post.Body[:147] + "..."
	} else {
		postShow.ShortDesc = post.Body
	}

	return &postShow
}
