package service

import (
	"bytes"
	"github.com/auyer/steganography"
	"image"
	"log"
)

type Steganography interface {
	EncodeTextToPic(image image.Image, text string) (*bytes.Buffer, error)
	DecodeTextToPic(image image.Image) (string, error)
	EncodePicToPic(image image.Image, image2Buf bytes.Buffer) (*bytes.Buffer, error)
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

func (s SteganographyService) DecodeTextToPic(image image.Image) (string, error) {
	sizeOfMessage := steganography.GetMessageSizeFromImage(image)
	msg := steganography.Decode(sizeOfMessage, image)
	return string(msg), nil
}
