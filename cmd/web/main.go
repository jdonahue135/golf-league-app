package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jdonahue135/golf-league-app/internal/config"
	"github.com/jdonahue135/golf-league-app/internal/driver"
	"github.com/jdonahue135/golf-league-app/internal/handlers"
	"github.com/jdonahue135/golf-league-app/internal/helpers"
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/render"
	"github.com/jdonahue135/golf-league-app/internal/repository/leaguerepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/playerrepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/userrepo"
	"github.com/jdonahue135/golf-league-app/internal/services/leagueservice"
	"github.com/jdonahue135/golf-league-app/internal/services/playerservice"
	"github.com/jdonahue135/golf-league-app/internal/services/userservice"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what will we put in the session
	gob.Register(models.User{})
	gob.Register(models.League{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database name")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	// change this to true when in production
	app.InProduction = *inProduction

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database.")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = *useCache

	userRepo := userrepo.NewPostgresUserRepo(db.SQL)
	userService := userservice.NewUserService(userRepo)
	playerRepo := playerrepo.NewPostgresPlayerRepo(db.SQL)
	playerService := playerservice.NewPlayerService(playerRepo)
	leagueRepo := leaguerepo.NewPostgresLeagueRepo(db.SQL)
	leagueService := leagueservice.NewLeagueService(leagueRepo, playerRepo, userRepo)
	handlers.NewHandlers(&app, userService, leagueService, playerService)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
