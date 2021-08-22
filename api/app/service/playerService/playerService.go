package playerService

import (
	"goa-golang/app/model/playerModel"
	"goa-golang/app/repository/playerRepository"
)

//PlayerServiceInterface define the player service interface methods
type PlayerServiceInterface interface {
	FindByID(uuid string) (player *playerModel.Player, err error)
}

// billingService handles communication with the player repository
type PlayerService struct {
	playerRepo playerRepository.PlayerRepositoryInterface
}

// NewPlayerService implements the player service interface.
func NewPlayerService(playerRepo playerRepository.PlayerRepositoryInterface) PlayerServiceInterface {
	return &PlayerService{
		playerRepo,
	}
}

// FindByID implements the method to find a player model by primary key
func (s *PlayerService) FindByID(uuid string) (player *playerModel.Player, err error) {
	return s.playerRepo.FindByID(uuid)
}
