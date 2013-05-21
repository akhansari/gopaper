package gopaper

import (
	"github.com/gorilla/mux"
	"net/http"

	"gopaper/handlers"
)

// $HOME/google_appengine/dev_appserver.py $HOME/Projects/gopaper/
// $HOME/google_appengine/appcfg.py update $HOME/Projects/gopaper/

func init() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", handlers.MakeHandler(handlers.Home.Index))
	router.HandleFunc("/page/{id:[0-9]+}", handlers.MakeHandler(handlers.Home.Index))
	router.HandleFunc("/post/{id:[0-9]+}/{title:[\\w-]+}", handlers.MakeHandler(handlers.Post.Index))

	router.HandleFunc("/tag/{tag:[\\w+-]+}", handlers.MakeHandler(handlers.Home.Tags))
	router.HandleFunc("/tag/{tag:[\\w+-]+}/page/{id:[0-9]+}", handlers.MakeHandler(handlers.Home.Tags))

	router.HandleFunc("/backend/", handlers.MakeHandler(handlers.Backend.Index))
	router.HandleFunc("/backend/page/{id:[0-9]+}", handlers.MakeHandler(handlers.Backend.Index))
	router.HandleFunc("/backend/addpost", handlers.MakeHandler(handlers.Backend.AddPost))
	router.HandleFunc("/backend/post/{id:[0-9]+}", handlers.MakeHandler(handlers.Backend.Post))
	router.HandleFunc("/backend/deletepost/{id:[0-9]+}", handlers.MakeHandler(handlers.Backend.DeletePost))

	http.Handle("/", router)
}
