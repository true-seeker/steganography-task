package main

import (
	"net/http"
	"steganography-task/internal/app"
)

func main() {
	r := http.NewServeMux()

	a := app.New(r)

	a.Run()
}
