package app

import (
	"net/http"
	"steganography-task/internal/router"
)

type App struct {
	r *http.ServeMux
}

func New(r *http.ServeMux) *App {
	return &App{r: router.New().InitRoutes(r)}
}

func (a *App) Run() error {
	err := http.ListenAndServe(":80", a.r)
	if err != nil {
		return err
	}
	return nil
}
