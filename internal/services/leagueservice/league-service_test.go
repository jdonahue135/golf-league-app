package leagueservice

import (
	"fmt"
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/models"
)

func TestGetLeague(t *testing.T) {
	service.GetLeague(1)
}

func TestGetLeagueByName(t *testing.T) {
	service.GetLeagueByName("name")
}

func TestGetLeaguesByUser(t *testing.T) {
	service.GetLeaguesByUser(1)
}

var createLeagueTests = []struct {
	name             string
	league           models.League
	commissioner     models.Player
	expectedLeagueID int
	expectError      bool
}{
	{
		"league error",
		models.League{Name: "League Error"},
		models.Player{},
		0,
		true,
	},
	{
		"player error",
		models.League{},
		models.Player{Handicap: 100},
		0,
		true,
	},
	{
		"success",
		models.League{},
		models.Player{},
		1,
		false,
	},
}

func TestCreateLeagueWithCommissioner(t *testing.T) {
	for _, e := range createLeagueTests {
		leagueID, err := service.CreateLeagueWithCommissioner(e.league, e.commissioner)
		if leagueID != e.expectedLeagueID {
			t.Errorf("failed %s: leagueID %d, but got %d", e.name, e.expectedLeagueID, leagueID)
		}
		if e.expectError && err == nil {
			t.Errorf("failed %s: expected error but got none", e.name)
		}
		if !e.expectError && err != nil {
			t.Errorf("failed %s: expected no error but got one", e.name)
		}
	}
}

var existingUserTests = []struct {
	name        string
	userID      int
	LeagueID    int
	expectError bool
}{
	{
		"error - player already active in league",
		3,
		0,
		true,
	},
	{
		"error - player inactive and can't activate",
		2,
		1,
		true,
	},
	{
		"error - create player db error",
		0,
		1,
		true,
	},
	{
		"success - existing player activate",
		4,
		2,
		false,
	},
	{
		"success - new player",
		0,
		2,
		false,
	},
}

func TestAddExistingUserToLeague(t *testing.T) {
	for _, e := range existingUserTests {
		err := service.AddExistingUserToLeague(e.userID, e.LeagueID)
		if e.expectError && err == nil {
			t.Errorf("failed %s: expected error but got none", e.name)
		}
		if !e.expectError && err != nil {
			fmt.Println(err)
			t.Errorf("failed %s: expected no error but got one", e.name)
		}
	}
}

var newUserTests = []struct {
	name        string
	user        models.User
	LeagueID    int
	expectError bool
}{
	{
		"error - user create error",
		models.User{FirstName: "user create error"},
		0,
		true,
	},
	{
		"error - player create error",
		models.User{FirstName: "player create error"},
		0,
		true,
	},
	{
		"success",
		models.User{},
		1,
		false,
	},
}

func TestAddNewUserToLeague(t *testing.T) {
	for _, e := range newUserTests {
		err := service.AddNewUserToLeague(e.user, e.LeagueID)
		if e.expectError && err == nil {
			t.Errorf("failed %s: expected error but got none", e.name)
		}
		if !e.expectError && err != nil {
			t.Errorf("failed %s: expected no error but got one", e.name)
		}
	}
}
