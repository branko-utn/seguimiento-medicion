package server

import "net/http"

type app struct{}

func nuevaApp(svc Servicios) (*app, error) { return &app{}, nil }

func (a *app) render(w http.ResponseWriter, r *http.Request, pagina string, datos any) {}

func (a *app) renderError(w http.ResponseWriter, r *http.Request, status int, mensaje string) {}
