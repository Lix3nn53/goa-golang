package characterRepository

import (
	"database/sql"
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
		var read1 sql.NullString
		var read2 sql.NullString
		var read3 sql.NullString
		var read4 sql.NullString
		var read5 sql.NullString
		var read6 sql.NullString
		var read7 sql.NullString
		var read8 sql.NullString
		if err := rows.Scan(&read1, &read2, &read3, &read4, &read5, &read6, &read7, &read8); err != nil {
			return nil, err
		}

		var characterNo string
		if read1.Valid {
			characterNo = read1.String
		}
		var chatTag string
		if read2.Valid {
			chatTag = read2.String
		}
		var craftingExperiences string
		if read3.Valid {
			craftingExperiences = read3.String
		}
		var turnedInQuests string
		if read4.Valid {
			turnedInQuests = read4.String
		}
		var activeQuests string
		if read1.Valid {
			activeQuests = read5.String
		}
		var rpgClass string
		if read2.Valid {
			rpgClass = read6.String
		}
		var unlockedClasses string
		if read3.Valid {
			unlockedClasses = read7.String
		}
		var totalExp string
		if read4.Valid {
			totalExp = read8.String
		}

		character := &characterModel.Character{
			CharacterNo:         characterNo,
			ChatTag:             chatTag,
			CraftingExperiences: craftingExperiences,
			TurnedInQuests:      turnedInQuests,
			ActiveQuests:        activeQuests,
			RpgClass:            rpgClass,
			UnlockedClasses:     unlockedClasses,
			TotalExp:            totalExp,
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
