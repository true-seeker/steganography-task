package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
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
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, "text/encode", http.StatusPermanentRedirect)
}

func (h *Handler) TextEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/templates/text_to_pic_encode.html")
}

func (h *Handler) TextDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/templates/text_to_pic_decode.html")
}

func (h *Handler) PictureEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/templates/pic_to_pic_encode.html")
}

func (h *Handler) PictureDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "ui/templates/pic_to_pic_decode.html")
}

func (h *Handler) Audio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "ui/templates/audio.html")
}

func (h *Handler) TextToPicEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inputTextForm := r.PostForm.Get("inputText")
	hostFileForm := r.PostForm.Get("hostFile")
	fmt.Println("inputTextForm : ", inputTextForm)
	fmt.Println("hostFileForm : ", len(hostFileForm))

	hostFile, hostFileHeader, hostFileErr := r.FormFile("hostFile")
	if hostFileErr != nil {
		http.Error(w, "hostFile is missing", http.StatusBadRequest)
		return
	}
	defer hostFile.Close()

	hostImage, _, err := image.Decode(hostFile)
	if err != nil {
		log.Printf("Host image decode error error %v", err)
		http.Error(w, "Host image decode error error", http.StatusBadRequest)
		return
	}

	encodedImageBuf, err := h.steganographyService.EncodeTextToPic(hostImage, inputTextForm)
	if err != nil {
		log.Printf("Eecode error %v", err)
		http.Error(w, "Encode error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", hostFileHeader.Filename))
	w.Header().Set("Content-Type", http.DetectContentType(encodedImageBuf.Bytes()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", encodedImageBuf.Len()))

	_, err = w.Write(encodedImageBuf.Bytes())
	if err != nil {
		log.Printf("Buf write error %v", err)
		http.Error(w, "Buf write error", http.StatusBadRequest)
		return
	}

	return
}

func (h *Handler) TextToPicDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//sourceFileForm := r.PostForm.Get("sourceFile")
	//fmt.Println("sourceFileForm : ", len(sourceFileForm))

	sourceFile, _, sourceFileErr := r.FormFile("sourceFile")
	if sourceFileErr != nil {
		http.Error(w, "sourceFile is missing", http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	sourceImage, _, err := image.Decode(sourceFile)
	if err != nil {
		log.Printf("Source image decode error error %v", err)
		http.Error(w, "Source image decode error error", http.StatusBadRequest)
		return
	}

	message, err := h.steganographyService.DecodeTextToPic(sourceImage)
	if err != nil {
		log.Printf("Decode error %v", err)
		http.Error(w, "Decode error", http.StatusBadRequest)
		return
	}
	fmt.Println(message)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		log.Printf("json encoder error %v", err)
		http.Error(w, "json encoder error", http.StatusBadRequest)
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

	sourceFile, _, sourceFileErr := r.FormFile("sourceFile")
	if sourceFileErr != nil {
		http.Error(w, "sourceFile is missing", http.StatusBadRequest)
		return
	}
	defer sourceFile.Close()

	hostImage, _, err := image.Decode(hostFile)
	if err != nil {
		log.Printf("Host image decode error %v", err)
		http.Error(w, "Host image decode error", http.StatusBadRequest)
		return
	}

	sourceImageBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(sourceImageBuf, hostFile); err != nil {
		log.Printf("Source image buffer error %v", err)
		http.Error(w, "Source image buffer error", http.StatusBadRequest)
		return
	}

	encodedImageBuf, err := h.steganographyService.EncodePicToPic(hostImage, *sourceImageBuf)
	if err != nil {
		log.Printf("Eecode error %v", err)
		http.Error(w, "Encode error", http.StatusBadRequest)
		return
	}

	fmt.Println(encodedImageBuf.Len())

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", hostFileHeader.Filename))
	w.Header().Set("Content-Type", http.DetectContentType(encodedImageBuf.Bytes()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", encodedImageBuf.Len()))

	_, err = w.Write(encodedImageBuf.Bytes())
	if err != nil {
		log.Printf("Buf write error %v", err)
		http.Error(w, "Buf write error", http.StatusBadRequest)
		return
	}
}
