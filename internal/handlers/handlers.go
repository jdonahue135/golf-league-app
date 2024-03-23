package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jdonahue135/golf-league-app/internal/config"
	"github.com/jdonahue135/golf-league-app/internal/forms"
	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/render"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

const leagueIDIndex = 2

var App *config.AppConfig

var Handler *Handlers

var UserService services.UserService

var LeagueService services.LeagueService

var PlayerService services.PlayerService

type Handlers struct {
	App           *config.AppConfig
	UserService   services.UserService
	LeagueService services.LeagueService
	PlayerService services.PlayerService
}

// NewHandlers sets dependencies of handlers
func NewHandlers(
	a *config.AppConfig,
	userService services.UserService,
	leagueService services.LeagueService,
	playerService services.PlayerService,
) {
	h := Handlers{
		App:           a,
		UserService:   userService,
		LeagueService: leagueService,
		PlayerService: playerService,
	}
	Handler = &h
}

// Home is the home page handler
func (m *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Handlers) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Leagues is the league page handler
func (m *Handlers) Leagues(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	userID, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "must be logged in to do that!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	_, err := m.UserService.GetUser(userID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "user not found!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	leagues, err := m.LeagueService.GetLeaguesByUser(userID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err.Error())
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["leagues"] = leagues

	render.Template(w, r, "leagues.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// ShowLeagueForm renders the create a league page and displays form
func (m *Handlers) ShowLeagueForm(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "create-league.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// ShowLeague shows information for a specific league
func (m *Handlers) ShowLeague(w http.ResponseWriter, r *http.Request) {
	leagueID, err := getLeagueIDFromURI(r.RequestURI)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	league, err := m.LeagueService.GetLeague(leagueID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find league")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	players, err := m.PlayerService.GetPlayersInLeague(league.ID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot get players for league")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data["league"] = league
	data["players"] = players

	render.Template(w, r, "league.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func getLeagueIDFromURI(URI string) (int, error) {
	exploded := strings.Split(URI, "/")
	return strconv.Atoi(exploded[leagueIDIndex])
}

// CreateLeague handles request to create a league
func (m *Handlers) CreateLeague(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "must be logged in to do that!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	_, err := m.UserService.GetUser(userID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "user not found!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	league := models.League{
		Name: r.Form.Get("name"),
	}

	commissioner := models.Player{
		UserID:         userID,
		IsCommissioner: true,
		IsActive:       true,
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
	_, err = m.LeagueService.GetLeagueByName(league.Name)

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
	id, err := m.LeagueService.CreateLeagueWithCommissioner(league, commissioner)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert league into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//redirect to league view page
	m.App.Session.Put(r.Context(), "league", league)
	http.Redirect(w, r, fmt.Sprintf("leagues/%d", id), http.StatusSeeOther)
}

// ShowAddPlayerForm renders the add player to a league page and displays form
func (m *Handlers) ShowAddPlayerForm(w http.ResponseWriter, r *http.Request) {
	leagueID, err := getLeagueIDFromURI(r.RequestURI)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})

	league, err := m.LeagueService.GetLeague(leagueID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find league")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data["league"] = league
	// send the data to the template
	render.Template(w, r, "add-player.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Handlers) AddPlayer(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "must be logged in to do that!")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	leagueID, err := getLeagueIDFromURI(r.RequestURI)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 2)
	form.MaxLength("first_name", 35)
	form.MinLength("last_name", 2)
	form.MaxLength("last_name", 35)
	form.IsEmail("email")

	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	email := r.Form.Get("email")
	playerUser := models.User{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		AccessLevel: models.AccessLevelPlayer,
	}
	if !form.Valid() {
		data := make(map[string]interface{})
		data["user"] = playerUser

		render.Template(w, r, "add-player.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	_, err = m.LeagueService.GetLeague(leagueID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot find league")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	player, err := m.PlayerService.GetPlayerInLeague(userID, leagueID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "you must be a member of this league to do that!")
		http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
		return
	}

	if !player.IsCommissioner {
		m.App.Session.Put(r.Context(), "error", "user must be league commissioner to add players!")
		http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
		return
	}

	playerUser, err = m.UserService.GetUserByEmail(email)
	if err == nil {
		//user already exists
		err = m.LeagueService.AddExistingUserToLeague(playerUser.ID, leagueID)
		if err != nil {
			m.App.Session.Put(r.Context(), "error", err.Error())
			http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
			return
		}
		m.App.Session.Put(r.Context(), "flash", "player added!")
		http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
		return
	}

	//user does not exist, need to create user and player records at same time
	err = m.LeagueService.AddNewUserToLeague(playerUser, leagueID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "error adding player to DB")
		http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "player added!")
	http.Redirect(w, r, fmt.Sprintf("/leagues/%d", leagueID), http.StatusSeeOther)
	return
}

// ShowSignUp shows the sign up page
func (m *Handlers) ShowSignUp(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "sign-up.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowSignUp handles logging the user in
func (m *Handlers) PostShowSignUp(w http.ResponseWriter, r *http.Request) {
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
	_, err = m.UserService.GetUserByEmail(email)
	if err == nil {
		form.Errors.Add("email", "Account already exists with that email address")
	}
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
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
	id, err := m.UserService.CreateUser(user, password)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert user into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", models.AccessLevelPlayer)
	m.App.Session.Put(r.Context(), "flash", "Signed up successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ShowLogin shows the login page
func (m *Handlers) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Logout logs the user out
func (m *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// PostShowLogin handles logging the user in
func (m *Handlers) PostShowLogin(w http.ResponseWriter, r *http.Request) {
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

	id, accessLevel, err := m.UserService.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "access_level", accessLevel)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Handlers) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}
