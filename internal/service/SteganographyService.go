package service

import (
	"bytes"
	"fmt"
	"github.com/auyer/steganography"
	"image"
	"log"
)

type Steganography interface {
	EncodeTextToPic(image image.Image, text string) (*bytes.Buffer, error)
	EncodePicToPic(image image.Image, image2Buf bytes.Buffer) (*bytes.Buffer, error)
	Decode(image image.Image) (string, error)
}

type SteganographyService struct {
}

func NewSteganographyService() Steganography {
	return &SteganographyService{}
}

func (s SteganographyService) EncodeTextToPic(image image.Image, text string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := steganography.Encode(buf, image, []byte(text))
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		return nil, err
	}
	return buf, nil
}

func (s SteganographyService) EncodePicToPic(image image.Image, image2Buf bytes.Buffer) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := steganography.Encode(buf, image, image2Buf.Bytes())
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		return nil, err
	}
	return buf, nil
}

func (s SteganographyService) Decode(image image.Image) (string, error) {
	sizeOfMessage := steganography.GetMessageSizeFromImage(image)

	msg := steganography.Decode(sizeOfMessage, image)
	fmt.Println(string(msg))
	return string(msg), nil
}
