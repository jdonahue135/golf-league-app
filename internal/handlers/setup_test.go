package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jdonahue135/golf-league-app/internal/config"
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/render"
	"github.com/jdonahue135/golf-league-app/internal/repository/leaguerepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/playerrepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/userrepo"
	"github.com/jdonahue135/golf-league-app/internal/services/leagueservice"
	"github.com/jdonahue135/golf-league-app/internal/services/playerservice"
	"github.com/jdonahue135/golf-league-app/internal/services/userservice"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

func TestMain(m *testing.M) {
	gob.Register(models.User{})
	gob.Register(map[string]int{})

	// change this to true when in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	userRepo := userrepo.NewTestUserRepo()
	userService := userservice.NewTestUserService(userRepo)
	playerRepo := playerrepo.NewTestPlayerRepo()
	playerService := playerservice.NewTestPlayerService(playerRepo)

	leagueRepo := leaguerepo.NewTestLeagueRepo()
	leagueService := leagueservice.NewTestLeagueService(leagueRepo, playerRepo, userRepo)
	NewHandlers(&app, userService, leagueService, playerService)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Handler.Home)
	mux.Get("/about", Handler.About)

	mux.Route("/leagues", func(mux chi.Router) {
		mux.Get("/", Handler.Leagues)
		mux.Post("/", Handler.CreateLeague)
		mux.Get("/new", Handler.ShowLeagueForm)
		mux.Get("/{id}", Handler.ShowLeague)
		mux.Get("/{id}/add-player", Handler.ShowAddPlayerForm)
		mux.Post("/{id}/players", Handler.AddPlayer)
	})

	mux.Route("/user", func(mux chi.Router) {
		mux.Get("/login", Handler.ShowLogin)
		mux.Post("/login", Handler.PostShowLogin)
		mux.Get("/logout", Handler.Logout)
		mux.Get("/sign-up", Handler.ShowSignUp)
		mux.Post("/sign-up", Handler.PostShowSignUp)
	})

	mux.Route("/admin", func(mux chi.Router) {
		mux.Get("/dashboard", Handler.AdminDashboard)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		log.Println(err)
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			log.Println(err)
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			log.Println(err)
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				log.Println(err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
