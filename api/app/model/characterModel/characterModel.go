package characterModel

import "database/sql"

// User represents user resources.
type Character struct {
	CharacterNo         string `json:"character_no" db:"character_no"`
	ChatTag             string `json:"chat_tag" db:"chat_tag"`
	CraftingExperiences string `json:"crafting_experiences" db:"crafting_experiences"`
	TurnedInQuests      string `json:"turnedinquests" db:"turnedinquests"`
	ActiveQuests        string `json:"activequests" db:"activequests"`
	RpgClass            string `json:"rpg_class" db:"rpg_class"`
	UnlockedClasses     string `json:"unlocked_classes" db:"unlocked_classes"`
	TotalExp            string `json:"totalexp" db:"totalexp"`
}

type CharacterScan struct {
	CharacterNo         sql.NullString `json:"character_no" db:"character_no"`
	ChatTag             sql.NullString `json:"chat_tag" db:"chat_tag"`
	CraftingExperiences sql.NullString `json:"crafting_experiences" db:"crafting_experiences"`
	TurnedInQuests      sql.NullString `json:"turnedinquests" db:"turnedinquests"`
	ActiveQuests        sql.NullString `json:"activequests" db:"activequests"`
	RpgClass            sql.NullString `json:"rpg_class" db:"rpg_class"`
	UnlockedClasses     sql.NullString `json:"unlocked_classes" db:"unlocked_classes"`
	TotalExp            sql.NullString `json:"totalexp" db:"totalexp"`
}
