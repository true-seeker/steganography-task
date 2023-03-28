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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	inputTextForm := r.PostForm.Get("inputText")

	hostFile, hostFileHeader, hostFileErr := r.FormFile("hostFile")
	if hostFileErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": hostFileErr.Error()})
		return
	}
	defer hostFile.Close()

	hostImage, _, err := image.Decode(hostFile)
	if err != nil {
		log.Printf("Host image decode error error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	encodedImageBuf, err := h.steganographyService.EncodeTextToPic(hostImage, inputTextForm)
	if err != nil {
		log.Printf("Encode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", hostFileHeader.Filename))
	w.Header().Set("Content-Type", http.DetectContentType(encodedImageBuf.Bytes()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", encodedImageBuf.Len()))

	_, err = w.Write(encodedImageBuf.Bytes())
	if err != nil {
		log.Printf("Buf write error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	sourceFile, _, sourceFileErr := r.FormFile("sourceFile")
	if sourceFileErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer sourceFile.Close()

	sourceImage, _, err := image.Decode(sourceFile)
	if err != nil {
		log.Printf("Source image decode error error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	message, err := h.steganographyService.DecodeTextToPic(sourceImage)
	if err != nil {
		log.Printf("Decode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"message": message})
	if err != nil {
		log.Printf("json encoder error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	return
}

func (h *Handler) PicToPicEncode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	hostFile, hostFileHeader, hostFileErr := r.FormFile("hostFile")
	if hostFileErr != nil {
		log.Printf("hostFile is missing %v", hostFileErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": hostFileErr.Error()})
		return
	}
	defer hostFile.Close()
	sourceFile, _, sourceFileErr := r.FormFile("sourceFile")
	if sourceFileErr != nil {
		log.Printf("sourceFile is missing %v", sourceFileErr)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": sourceFileErr.Error()})
		return
	}
	defer sourceFile.Close()

	hostImage, _, err := image.Decode(hostFile)
	if err != nil {
		log.Printf("Host image decode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	sourceImageBuf := bytes.NewBuffer(nil)
	if _, err := io.Copy(sourceImageBuf, sourceFile); err != nil {
		log.Printf("Source image buffer error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	encodedImageBuf, err := h.steganographyService.EncodePicToPic(hostImage, *sourceImageBuf)
	if err != nil {
		log.Printf("Encode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", hostFileHeader.Filename))
	w.Header().Set("Content-Type", http.DetectContentType(encodedImageBuf.Bytes()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", encodedImageBuf.Len()))

	_, err = w.Write(encodedImageBuf.Bytes())
	if err != nil {
		log.Printf("Buf write error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
}

func (h *Handler) PicToPicDecode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	hostFile, sourceFileHeader, sourceFileErr := r.FormFile("hostFile")
	if sourceFileErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": sourceFileErr.Error()})
		return
	}
	defer hostFile.Close()

	hostImage, _, err := image.Decode(hostFile)
	if err != nil {
		log.Printf("Host image decode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	decodedPicBuf, err := h.steganographyService.DecodePicToPic(hostImage)
	if err != nil {
		log.Printf("Decode error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", sourceFileHeader.Filename))
	w.Header().Set("Content-Type", http.DetectContentType(decodedPicBuf.Bytes()))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", decodedPicBuf.Len()))

	_, err = w.Write(decodedPicBuf.Bytes())
	if err != nil {
		log.Printf("Buf write error %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

}
