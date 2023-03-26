package service

type Steganography interface {
}

type SteganographyService struct {
}

func NewSteganographyService() Steganography {
	return &SteganographyService{}
}
