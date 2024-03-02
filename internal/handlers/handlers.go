package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jdonahue135/golf-league-app/internal/config"
	"github.com/jdonahue135/golf-league-app/internal/driver"
	"github.com/jdonahue135/golf-league-app/internal/forms"
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/render"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/repository/dbrepo"
)

const leagueIDIndex = 2

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Leagues is the about page handler
func (m *Repository) Leagues(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "leagues.page.tmpl", &models.TemplateData{})
}

// League renders the create a league page and displays form
func (m *Repository) League(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "create-league.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// ShowLeague shows information for a specific league
func (m *Repository) ShowLeague(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	leagueID, err := strconv.Atoi(exploded[leagueIDIndex])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	league, ok := m.App.Session.Get(r.Context(), "league").(models.League)
	if ok {
		data["league"] = league
	} else {
		league, err = m.DB.GetLeagueByID(leagueID)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", "cannot find league")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		data["league"] = league
	}

	render.Template(w, r, "league.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// CreateLeague handles request to create a league
func (m *Repository) CreateLeague(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	league := models.League{
		Name: r.Form.Get("name"),
	}

	form := forms.New(r.PostForm)

	form.Required("name")
	form.MinLength("name", 3)
	form.MaxLength("name", 50)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["league"] = league

		render.Template(w, r, "create-league.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//check if name is unique in db
	_, err = m.DB.GetLeagueByName(league.Name)

	if err == nil {
		form.Errors.Add("name", "This league name is taken, please choose another")
		data := make(map[string]interface{})
		data["league"] = league

		render.Template(w, r, "create-league.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//insert into db
	id, err := m.DB.CreateLeague(league)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert league into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//redirect to league view page
	m.App.Session.Put(r.Context(), "league", league)
	http.Redirect(w, r, fmt.Sprintf("leagues/%d", id), http.StatusSeeOther)
}

// ShowSignUp shows the sign up page
func (m *Repository) ShowSignUp(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "sign-up.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowSignUp handles logging the user in
func (m *Repository) PostShowSignUp(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "password")
	form.MinLength("first_name", 2)
	form.MaxLength("first_name", 35)
	form.MinLength("last_name", 2)
	form.MaxLength("last_name", 35)
	form.MinLength("password", 2)
	form.MaxLength("password", 35)
	form.IsEmail("email")

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	_, err = m.DB.GetUserByEmail(email)
	if err == nil {
		form.Errors.Add("email", "Account already exists with that email address")
	}

	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = user

		render.Template(w, r, "sign-up.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	password := r.Form.Get("password")
	id, err := m.DB.CreateUser(user, password)

	if err != nil {
		fmt.Println(err)
		m.App.Session.Put(r.Context(), "error", "can't insert user into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Signed up successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowLogin shows the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Logout logs the user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
