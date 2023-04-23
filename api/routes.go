package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	Handler RedirectHandler
	Router  *chi.Mux
}

func NewServer(handler RedirectHandler, app *chi.Mux) *Server {

	app.Use(middleware.RequestID)
	app.Use(middleware.RealIP)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	//app.Use(middleware.Timeout(60 * time.Second))

	app.Get("/{code}", handler.Get)
	app.Post("/", handler.Post)
	app.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ping"))

	})

	return &Server{Router: app, Handler: handler}
}
