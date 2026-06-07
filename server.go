package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"todo/components"

	"github.com/a-h/templ"
)

type Server struct {
	mux *http.ServeMux
	db  *Database
}

func NewServer(db *Database) *Server {
	return &Server{
		mux: http.NewServeMux(),
		db:  db,
	}
}

func (s *Server) Run() {
	s.Routes()
	fmt.Println("Listening on :6969")
	log.Fatal(http.ListenAndServe(":6969", s.mux))
}

func (s *Server) Routes() {
	s.mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	s.mux.HandleFunc("GET /{$}", layout_page)
	s.mux.HandleFunc("GET /home/", home_page)
	s.mux.HandleFunc("GET /create/", create_page)
	s.mux.HandleFunc("POST /create-todo", s.db.create_todo)
	s.mux.HandleFunc("GET /show", s.db.show_page)
	s.mux.HandleFunc("GET /favicon.ico", load_favicon)
}

func RenderLayout(ctx context.Context, w io.Writer, component templ.Component) {
	layout := components.Layout("Greeting", component)
	layout.Render(ctx, w)
}

func layout_page(w http.ResponseWriter, r *http.Request) {
	home := components.Home()
	RenderLayout(r.Context(), w, home)
}

func home_page(w http.ResponseWriter, r *http.Request) {
	home := components.Home()

	if r.Header.Get("Hx-Request") == "true" {
		home.Render(r.Context(), w)
		return
	}

	RenderLayout(r.Context(), w, home)
}

func create_page(w http.ResponseWriter, r *http.Request) {
	create := components.Create()
	if r.Header.Get("Hx-Request") == "true" {
		create.Render(r.Context(), w)
		return
	}

	RenderLayout(r.Context(), w, create)
}

func load_favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
