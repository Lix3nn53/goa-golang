package characterRepository

import (
	"goa-golang/app/model/characterModel"
	"goa-golang/internal/storage"
)

// billingRepository handles communication with the user store
type CharacterRepository struct {
	db *storage.DbStore
}

//CharacterRepositoryInterface define the user repository interface methods
type CharacterRepositoryInterface interface {
	FindByID(uuid string) (user []*characterModel.Character, err error)
	RemoveByID(uuid string, characterNo int) error
}

// NewCharacterRepository implements the user repository interface.
func NewCharacterRepository(db *storage.DbStore) CharacterRepositoryInterface {
	return &CharacterRepository{
		db,
	}
}

// FindByID implements the method to find a user from the store
func (r *CharacterRepository) FindByID(uuid string) (characters []*characterModel.Character, err error) {
	var query = "SELECT character_no, chat_tag, crafting_experiences, turnedinquests, activequests, rpg_class, unlocked_classes, totalexp FROM goa_player_character WHERE uuid = ?"
	rows, err := r.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	characters = make([]*characterModel.Character, 0)
	for rows.Next() {
		scan := &characterModel.CharacterScan{}

		if err := rows.Scan(&scan.CharacterNo, &scan.ChatTag, &scan.CraftingExperiences, &scan.TurnedInQuests, &scan.ActiveQuests,
			&scan.RpgClass, &scan.UnlockedClasses, &scan.TotalExp); err != nil {
			return nil, err
		}

		character := &characterModel.Character{
			CharacterNo:         scan.CharacterNo.String,
			ChatTag:             scan.ChatTag.String,
			CraftingExperiences: scan.CraftingExperiences.String,
			TurnedInQuests:      scan.TurnedInQuests.String,
			ActiveQuests:        scan.ActiveQuests.String,
			RpgClass:            scan.RpgClass.String,
			UnlockedClasses:     scan.UnlockedClasses.String,
			TotalExp:            scan.TotalExp.String,
		}

		characters = append(characters, character)
	}

	return characters, nil
}

// RemoveByID implements the method to remove a user from the store
func (r *CharacterRepository) RemoveByID(uuid string, characterNo int) error {

	_, err := r.db.Exec(`DELETE FROM goa_player_character WHERE uuid = ? AND character_no = ?;`, uuid, characterNo)
	return err
}
