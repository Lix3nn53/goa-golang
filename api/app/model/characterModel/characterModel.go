package characterModel

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
