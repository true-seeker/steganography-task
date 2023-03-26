package handlers

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
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

	err := r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stegoTypeForm := r.PostForm.Get("stegoType")
	inputTextForm := r.PostForm.Get("inputText")
	sourceFileForm := r.PostForm.Get("sourceFile")
	hostFileForm := r.PostForm.Get("hostFile")
	fmt.Println("stegoTypeForm : ", stegoTypeForm)
	fmt.Println("inputTextForm : ", inputTextForm)
	fmt.Println("sourceFileForm : ", len(sourceFileForm))
	fmt.Println("hostFileForm : ", len(hostFileForm))

	hostFile, hostFileHeader, hostFileErr := r.FormFile("hostFile")
	if hostFileErr != nil {
		http.Error(w, "hostFile is missing", http.StatusBadRequest)
		return
	}
	defer hostFile.Close()

	out, err := os.Create(hostFileHeader.Filename)
	defer out.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.Copy(out, hostFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	return

}

func (h *Handler) PicToPic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}
	stegoTypeForm := r.PostForm.Get("stegoType")
	inputTextForm := r.PostForm.Get("inputText")
	sourceFileForm := r.PostForm.Get("sourceFile")
	hostFileForm := r.PostForm.Get("hostFile")
	fmt.Println("stegoTypeForm : ", stegoTypeForm)
	fmt.Println("inputTextForm : ", inputTextForm)
	fmt.Println("sourceFileForm : ", len(sourceFileForm))
	fmt.Println("hostFileForm : ", len(hostFileForm))

	hostFile, hostFileHeader, hostFileErr := r.FormFile("hostFile")
	if hostFileErr != nil {
		http.Error(w, "hostFile is missing", http.StatusBadRequest)
		return
	}
	defer hostFile.Close()

	sourceFile, sourceFileHeader, sourceFileErr := r.FormFile("sourceFile")
	if sourceFileErr != nil {
		http.Error(w, "sourceFile is missing", http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	out, err := os.Create(hostFileHeader.Filename)
	defer out.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.Copy(out, hostFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out2, err := os.Create(sourceFileHeader.Filename)
	defer out2.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = io.Copy(out2, sourceFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "PicToPic, %q", html.EscapeString(r.URL.Path))
}
