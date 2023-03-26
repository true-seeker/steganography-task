package handlers

import (
	"fmt"
	"html"
	"net/http"
	"steganography-task/internal/service"
)

type Handler struct {
	steganographyService service.Steganography
}

func NewHandler(steganographyService service.Steganography) *Handler {
	return &Handler{steganographyService: steganographyService}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Check path here
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/templates/index.html")
}

func (h *Handler) Audio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "ui/templates/audio.html")
}

func (h *Handler) TextToPic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "TextToPic, %q", html.EscapeString(r.URL.Path))
}

func (h *Handler) PicToPic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "PicToPic, %q", html.EscapeString(r.URL.Path))
}
