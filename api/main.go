package main

import (
	// "awesomeProject/api/session"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type MiddleWareHandler struct {
	r *httprouter.Router
}

func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := MiddleWareHandler{}
	m.r = r
	return m
}

func (m MiddleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/user", CreateUser)
	router.POST("/user/login", Login)
	router.GET("/user/:username", GetUserInfo)
	// router.POST("/user/:username/videos", AddNewVideo)
	// router.GET("/user/:username/videos", ListAllVideos)
	// router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)
	// router.POST("/videos:vid-id/comments", PostComment)
	// router.GET("/videos/:vid-id/comments", ShowComments)

	return router
}

// func Prepare() {
// 	session.LoadSessionFromDb()
// }

func main() {
	// Prepare()
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	http.ListenAndServe(":8000", mh)
}
