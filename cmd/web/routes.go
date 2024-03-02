package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jdonahue135/golf-league-app/internal/config"
	"github.com/jdonahue135/golf-league-app/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Route("/leagues", func(mux chi.Router) {
		mux.Get("/", handlers.Repo.Leagues)
		mux.Post("/", handlers.Repo.CreateLeague)
		mux.Get("/create-league", handlers.Repo.League)
		mux.Get("/{id}", handlers.Repo.ShowLeague)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/login", handlers.Repo.ShowLogin)
		mux.Post("/login", handlers.Repo.PostShowLogin)
		mux.Get("/logout", handlers.Repo.Logout)
		mux.Get("/sign-up", handlers.Repo.ShowSignUp)
		mux.Post("/sign-up", handlers.Repo.PostShowSignUp)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
