-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_character` (
     `character_no` smallint NOT NULL,
     `uuid` varchar(40) NOT NULL,
     `off_hand` text NULL,
     `slot_parrot` text NULL,
     `slot_necklace` text NULL,
     `slot_ring` text NULL,
     `slot_earring` text NULL,
     `slot_glove` text NULL,
     `slot_pet` text NULL,
     `chat_tag` varchar(45) NULL,
     `crafting_experiences` text NOT NULL,
     `inventory` mediumtext NOT NULL,
     `turnedinquests` text NULL,
     `activequests` text NULL,
     `location` text NOT NULL,
     `armor_content` text NOT NULL,
     `rpg_class` varchar(45) NOT NULL,
     `totalexp` int NOT NULL,
     `slot_tool_axe` text NOT NULL,
     `slot_tool_bottle` text NOT NULL,
     `slot_tool_hoe` text NOT NULL,
     `slot_tool_pickaxe` text NOT NULL,
     UNIQUE KEY `Ind_88` (`uuid`, `character_no`),
     KEY `fkIdx_55` (`uuid`),
     CONSTRAINT `FK_55` FOREIGN KEY `fkIdx_55` (`uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE goa_player_character;