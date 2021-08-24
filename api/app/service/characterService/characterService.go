package characterService

import (
	"goa-golang/app/model/characterModel"
	"goa-golang/app/repository/characterRepository"
)

//CharacterServiceInterface define the character service interface methods
type CharacterServiceInterface interface {
	FindByID(uuid string) (character []*characterModel.Character, err error)
}

// billingService handles communication with the character repository
type CharacterService struct {
	characterRepo characterRepository.CharacterRepositoryInterface
}

// NewCharacterService implements the character service interface.
func NewCharacterService(characterRepo characterRepository.CharacterRepositoryInterface) CharacterServiceInterface {
	return &CharacterService{
		characterRepo,
	}
}

// FindByID implements the method to find a character model by primary key
func (s *CharacterService) FindByID(uuid string) (character []*characterModel.Character, err error) {
	return s.characterRepo.FindByID(uuid)
}
