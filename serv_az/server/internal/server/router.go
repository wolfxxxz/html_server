package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	GetPost(string, http.HandlerFunc)
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
	InitImages(string, string)
}

type router struct {
	mux *mux.Router
}

func (router *router) InitImages(path string, absolutePath string) {
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("web/images"))))
	//srv.mux.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("templates/images"))))
	router.mux.PathPrefix(path).Handler(http.StripPrefix(path, http.FileServer(http.Dir(absolutePath))))
}

func (router *router) ServeHttp(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}

func (router *router) GetPost(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodGet, http.MethodPost)
}

func (router *router) Get(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodGet)
}

func (router *router) Post(path string, handlerFunc http.HandlerFunc) {
	router.mux.HandleFunc(path, handlerFunc).Methods(http.MethodPost)
}
