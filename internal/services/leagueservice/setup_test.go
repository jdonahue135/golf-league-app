package leagueservice

import (
	"os"
	"testing"

	"github.com/jdonahue135/golf-league-app/internal/repository/leaguerepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/playerrepo"
	"github.com/jdonahue135/golf-league-app/internal/repository/userrepo"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

var service services.LeagueService

func TestMain(m *testing.M) {
	leagueRepo := leaguerepo.NewTestLeagueRepo()
	playerRepo := playerrepo.NewTestPlayerRepo()
	userRepo := userrepo.NewTestUserRepo()
	service = NewLeagueService(leagueRepo, playerRepo, userRepo)

	os.Exit(m.Run())
}
