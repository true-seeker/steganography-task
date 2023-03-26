package router

import (
	"net/http"
	"steganography-task/internal/handlers"
	"steganography-task/internal/service"
)

type Router struct {
}

func New() *Router {
	return &Router{}
}

func (router *Router) InitRoutes(r *http.ServeMux) *http.ServeMux {
	steganographyService := service.NewSteganographyService()
	handler := handlers.NewHandler(steganographyService)

	r.HandleFunc("/", handler.Index)
	r.HandleFunc("/audio", handler.Audio)

	return r
}
