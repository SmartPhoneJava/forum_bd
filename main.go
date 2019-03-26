package main

import (
	"escapade/internal/services/api"
	"fmt"
	"os"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	const confPath = "conf.json"

	API, conf, err := api.GetHandler(confPath) // init.go
	if err != nil {
		fmt.Println("Some error happened with configuration file or database" + err.Error())
		return
	}

	r := mux.NewRouter()

	var api = r.PathPrefix("/api").Subrouter()
	var user = api.PathPrefix("/user").Subrouter()

	user.HandleFunc("/{name}/create", API.CreateUser).Methods("POST")
	user.HandleFunc("/{name}/profile", API.GetProfile).Methods("GET")
	user.HandleFunc("/{name}/profile", API.UpdateProfile).Methods("POST")

	var forum = api.PathPrefix("/forum").Subrouter()
	forum.HandleFunc("/create", API.CreateForum).Methods("POST")
	forum.HandleFunc("/{slug}/details", API.GetForum).Methods("GET")

	forum.HandleFunc("/{slug}/create", API.CreateThread).Methods("POST")
	forum.HandleFunc("/{slug}/threads", API.GetThreads).Methods("GET")

	forum.HandleFunc("/{slug}/users", API.GetUsers).Methods("GET")

	var thread = api.PathPrefix("/thread").Subrouter()
	thread.HandleFunc("/{slug}/details", API.GetThreadDetails).Methods("GET")
	thread.HandleFunc("/{slug}/details", API.UpdateThread).Methods("POST")
	thread.HandleFunc("/{slug}/create", API.CreatePosts).Methods("POST")
	thread.HandleFunc("/{slug}/posts", API.GetPosts).Methods("GET")

	thread.HandleFunc("/{slug}/vote", API.Vote).Methods("POST")

	var post = api.PathPrefix("/post").Subrouter()
	post.HandleFunc("/{id}/details", API.GetPostfull).Methods("GET")
	post.HandleFunc("/{id}/details", API.UpdatePost).Methods("POST")

	var service = api.PathPrefix("/service").Subrouter()
	service.HandleFunc("/status", API.GetStatus).Methods("GET")
	service.HandleFunc("/clear", API.Clear).Methods("POST")

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	fmt.Println("launched, look at us on " + conf.Server.Host + os.Getenv("PORT"))

	if err = http.ListenAndServe(os.Getenv("PORT"), r); err != nil { //os.Getenv("PORT"), r); err != nil {
		fmt.Println("oh, this is error:" + err.Error())
	}
}
