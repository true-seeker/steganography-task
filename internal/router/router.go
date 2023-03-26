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
	r.HandleFunc("/api/text_to_pic", handler.TextToPic)
	r.HandleFunc("/api/pic_to_pic", handler.PicToPic)

	return r
}
