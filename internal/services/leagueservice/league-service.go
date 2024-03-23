package leagueservice

import (
	"errors"

	"github.com/jdonahue135/golf-league-app/internal/models"
	"github.com/jdonahue135/golf-league-app/internal/repository"
	"github.com/jdonahue135/golf-league-app/internal/services"
)

type leagueService struct {
	LeagueRepo repository.LeagueRepo
	PlayerRepo repository.PlayerRepo
	UserRepo   repository.UserRepo
}

func NewLeagueService(l repository.LeagueRepo, p repository.PlayerRepo, u repository.UserRepo) services.LeagueService {
	return &leagueService{LeagueRepo: l, PlayerRepo: p, UserRepo: u}
}

func (m *leagueService) GetLeague(ID int) (models.League, error) {
	return m.LeagueRepo.GetLeagueByID(ID)
}

func (m *leagueService) GetLeagueByName(name string) (models.League, error) {
	return m.LeagueRepo.GetLeagueByName(name)
}

func (m *leagueService) GetLeaguesByUser(userID int) ([]models.League, error) {
	return m.LeagueRepo.GetLeaguesByUserID(userID)
}

func (m *leagueService) CreateLeagueWithCommissioner(league models.League, commissioner models.Player) (int, error) {
	ctx, cancel, tx, err := m.LeagueRepo.BeginTransaction()
	defer cancel()
	if err != nil {
		return 0, err
	}

	leagueID, err := m.LeagueRepo.CreateLeagueTransaction(league, ctx, tx)
	if err != nil {
		return 0, err
	}

	commissioner.LeagueID = leagueID
	err = m.PlayerRepo.CreatePlayerTransaction(commissioner, ctx, tx)
	if err != nil {
		return 0, err
	}

	err = m.PlayerRepo.CommitTransaction(tx)
	if err != nil {
		return 0, err
	}

	return leagueID, nil
}

func (m *leagueService) AddExistingUserToLeague(userID, leagueID int) error {
	player, err := m.PlayerRepo.GetPlayerByUserAndLeagueID(userID, leagueID)
	if err == nil {
		if player.IsActive {
			return errors.New("this player is already in this league")
		}
		player.IsActive = true
		err = m.PlayerRepo.UpdatePlayer(player)
		if err != nil {
			return errors.New("cannot reactivate player")
		}
		return nil
	}

	player = models.Player{
		LeagueID:       leagueID,
		UserID:         userID,
		IsActive:       true,
		IsCommissioner: false,
	}
	err = m.PlayerRepo.CreatePlayer(player)
	return err
}

func (m *leagueService) AddNewUserToLeague(user models.User, leagueID int) error {
	//transaction
	ctx, cancel, tx, err := m.UserRepo.BeginTransaction()
	defer cancel()
	if err != nil {
		return err
	}

	//create user
	userID, err := m.UserRepo.CreateInactiveUserTransaction(user, ctx, tx)
	if err != nil {
		return err
	}

	//create player
	player := models.Player{
		UserID:         userID,
		LeagueID:       leagueID,
		IsCommissioner: false,
		IsActive:       true,
	}
	err = m.PlayerRepo.CreatePlayerTransaction(player, ctx, tx)
	if err != nil {
		return err
	}

	err = m.PlayerRepo.CommitTransaction(tx)

	if err != nil {
		return err
	}

	return nil
}
