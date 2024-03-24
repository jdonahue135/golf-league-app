package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
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
	{"non-existent", "/green/eggs/and/ham", "GET", http.StatusNotFound},
	{"about", "/about", "GET", http.StatusOK},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"sign up", "/user/sign-up", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"create league", "/leagues/new", "GET", http.StatusOK},
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
		handler := http.HandlerFunc(Handler.ShowLeague)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var leaguesTests = []struct {
	name               string
	userID             int
	expectedStatusCode int
}{
	{
		name:               "no user logged in",
		userID:             -1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not found",
		userID:             0,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "league service error",
		userID:             2,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "success",
		userID:             1,
		expectedStatusCode: http.StatusOK,
	},
}

func TestLeagues(t *testing.T) {
	for _, e := range leaguesTests {
		req, _ := http.NewRequest("GET", "/leagues", nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		if e.userID >= 0 {
			session.Put(req.Context(), "user_id", e.userID)
		}

		handler := http.HandlerFunc(Handler.Leagues)
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

		handler := http.HandlerFunc(Handler.CreateLeague)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}

}

var showPlayerTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{
		name:               "bad url parameter",
		url:                "/leagues/s/add-player",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "non-existing league",
		url:                "/leagues/3/add-player",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "existing league",
		url:                "/leagues/1/add-player",
		expectedStatusCode: http.StatusOK,
	},
}

func TestShowAddPlayerForm(t *testing.T) {
	for _, e := range showPlayerTests {
		req, _ := http.NewRequest("GET", e.url, nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.RequestURI = e.url

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Handler.ShowAddPlayerForm)
		handler.ServeHTTP(rr, req)
		if rr.Code != e.expectedStatusCode {
			t.Errorf("%s returned wrong response code: got %d, wanted %d", e.name, rr.Code, e.expectedStatusCode)
		}
	}
}

var postPlayerTests = []struct {
	name               string
	firstName          string
	lastName           string
	email              string
	userID             int
	leagueID           int
	expectedStatusCode int
}{
	{
		name:               "user not logged in",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             -1,
		leagueID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not found",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             0,
		leagueID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not a member of league",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             4,
		leagueID:           4,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not commissioner of league",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             3,
		leagueID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "first name too short",
		firstName:          "J",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "first name too long",
		firstName:          "league0league0league0league0league0league0league0league0",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "last name too short",
		firstName:          "John",
		lastName:           "D",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "last name too long",
		firstName:          "John",
		lastName:           "league0league0league0league0league0league0league0league0",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "invalid email",
		firstName:          "John",
		lastName:           "Doe",
		email:              "johncom",
		userID:             1,
		leagueID:           1,
		expectedStatusCode: http.StatusOK,
	},
	{
		name:               "league doesn't exist",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           3,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "player not commissioner",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             2,
		leagueID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user exists - active in league",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             6,
		leagueID:           6,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user does not exist - error",
		firstName:          "John",
		lastName:           "Doe",
		email:              "me@here.ca",
		userID:             10,
		leagueID:           2,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user does not exist - success",
		firstName:          "John",
		lastName:           "Doe",
		email:              "me@here.ca",
		userID:             10,
		leagueID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "invalid url param",
		firstName:          "John",
		lastName:           "Doe",
		email:              "john@doe.com",
		userID:             1,
		leagueID:           0,
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestAddPlayer(t *testing.T) {
	for _, e := range postPlayerTests {
		postedData := url.Values{}
		postedData.Add("first_name", e.firstName)
		postedData.Add("last_name", e.lastName)
		postedData.Add("email", e.email)

		URI := fmt.Sprintf("/leagues/%d/players", e.leagueID)
		if e.leagueID == 0 {
			URI = "/leagues/s/players"
		}

		req, _ := http.NewRequest("POST", URI, strings.NewReader(postedData.Encode()))
		req.RequestURI = URI

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		if e.userID >= 0 {
			session.Put(req.Context(), "user_id", e.userID)
		}

		handler := http.HandlerFunc(Handler.AddPlayer)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

var deletePlayerTests = []struct {
	name               string
	userID             int
	leagueID           int
	playerID           int
	expectedStatusCode int
}{
	{
		name:               "user not logged in",
		userID:             -1,
		leagueID:           1,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "invalid league url param",
		userID:             1,
		leagueID:           0,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "invalid player url param",
		userID:             1,
		leagueID:           2,
		playerID:           0,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not found in league",
		userID:             4,
		leagueID:           4,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "user not commissioner in league",
		userID:             3,
		leagueID:           4,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "league doesn't exist",
		userID:             1,
		leagueID:           3,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "player doesn't exist",
		userID:             1,
		leagueID:           2,
		playerID:           9,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "player inactive in league",
		userID:             1,
		leagueID:           2,
		playerID:           8,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "service error",
		userID:             1,
		leagueID:           2,
		playerID:           10,
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "success",
		userID:             1,
		leagueID:           2,
		playerID:           1,
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestRemovePlayer(t *testing.T) {
	for _, e := range deletePlayerTests {

		URI := fmt.Sprintf("/leagues/%d/players/%d", e.leagueID, e.playerID)
		if e.leagueID == 0 {
			URI = "/leagues/s/players/1"
		}
		if e.playerID == 0 {
			URI = "/leagues/2/players/s"
		}

		req, _ := http.NewRequest("DELETE", URI, nil)
		req.RequestURI = URI

		ctx := getCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		if e.userID >= 0 {
			session.Put(req.Context(), "user_id", e.userID)
		}

		handler := http.HandlerFunc(Handler.RemovePlayer)
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
		handler := http.HandlerFunc(Handler.PostShowSignUp)
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

func TestPostShowLogin(t *testing.T) {
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
		handler := http.HandlerFunc(Handler.PostShowLogin)
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
