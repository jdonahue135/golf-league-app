package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/driver"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"sign up", "/user/sign-up", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"create league", "/leagues/create-league", "GET", http.StatusOK},
	{"leagues", "/leagues", "GET", http.StatusOK},
}

// TestHandlers tests all routes that don't require extra tests (gets)
func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

var leagueTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "existing league",
		url:                "/leagues/1",
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "non-existing league",
		url:                "/leagues/3",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "bad url parameter",
		url:                "/leagues/s",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "league with player error",
		url:                "/leagues/2",
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestShowLeague(t *testing.T) {
	for _, e := range leagueTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the RequestURI on the request so that we can grab the ID from the URL
		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.ShowLeague)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var postLeagueTests = []struct {
	name               string
	leagueName         string
	userID             int
	expectedStatusCode int
}{
	{
		name:               "user not logged in",
		leagueName:         "league0",
		userID:             -1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not found",
		leagueName:         "league0",
		userID:             0,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "name too short",
		leagueName:         "l",
		userID:             1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "name too long",
		leagueName:         "league0league0league0league0league0league0league0league0",
		userID:             1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "name not unique in DB",
		leagueName:         "league0",
		userID:             1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "error inserting league",
		leagueName:         "league1",
		userID:             1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "happy path",
		leagueName:         "league2",
		userID:             1,
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestCreateLeague(t *testing.T) {
	for _, e := range postLeagueTests {
		postedData := url.Values{}
		postedData.Add("name", e.leagueName)

		req, _ := http.NewRequest("POST", "/leagues", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		if e.userID >= 0 {
			session.Put(req.Context(), "user_id", e.userID)
		}

		handler := http.HandlerFunc(Repo.CreateLeague)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}

}

// loginTests is the data for the Login handler tests
var signUpTests = []struct {
	name               string
	firstName          string
	lastName           string
	password           string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid",
		"Jake",
		"Donahue",
		"password",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-data",
		"Jake",
		"Donahue",
		"password",
		"j",
		http.StatusOK,
		`action="/user/sign-up"`,
		"",
	},
	{
		"create user error",
		"Jake",
		"Donahue",
		"error",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
}

func TestPostShowSignUp(t *testing.T) {
	// range through all tests
	for _, e := range signUpTests {
		postedData := url.Values{}
		postedData.Add("first_name", e.firstName)
		postedData.Add("last_name", e.lastName)
		postedData.Add("email", e.email)
		postedData.Add("password", e.password)

		// create request
		req, _ := http.NewRequest("POST", "/user/sign-up", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowSignUp)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

// loginTests is the data for the Login handler tests
var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"me@here.ca",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jack@nimble.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {
	// range through all tests
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request
		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			// get the URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, actualLoc.String())
			}
		}

		// checking for expected values in HTML
		if e.expectedHTML != "" {
			// read the response body into a string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

// gets the context
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
