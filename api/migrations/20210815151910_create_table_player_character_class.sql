-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_character_class` (
     `class_name` NOT NULL,
     `uuid` varchar(40) NOT NULL,
     `character_no` smallint NOT NULL,
     `skill_points` mediumtext NOT NULL,
     `totalexp` int NOT NULL,
     `attribute_points` mediumtext NOT NULL,
     `skill_bar` mediumtext NOT NULL,
     UNIQUE KEY `Index_228` (`class_name`, `uuid`, `character_no`),
     KEY `FK_215` (`uuid`),
     CONSTRAINT `FK_213` FOREIGN KEY `FK_215` (`uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE goa_player_character_class;