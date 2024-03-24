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

	mux.Get("/", handlers.Handler.Home)
	mux.Get("/about", handlers.Handler.About)

	mux.Route("/leagues", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/", handlers.Handler.Leagues)
		mux.Post("/", handlers.Handler.CreateLeague)
		mux.Get("/new", handlers.Handler.ShowLeagueForm)
		mux.Get("/{id}", handlers.Handler.ShowLeague)
		mux.Get("/{id}/add-player", handlers.Handler.ShowAddPlayerForm)
		mux.Post("/{id}/players", handlers.Handler.AddPlayer)
		mux.Get("/{league_id}/players/{id}/remove-player", handlers.Handler.RemovePlayer)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/login", handlers.Handler.ShowLogin)
		mux.Post("/login", handlers.Handler.PostShowLogin)
		mux.Get("/logout", handlers.Handler.Logout)
		mux.Get("/sign-up", handlers.Handler.ShowSignUp)
		mux.Post("/sign-up", handlers.Handler.PostShowSignUp)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(AuthAdmin)

		mux.Get("/dashboard", handlers.Handler.AdminDashboard)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
