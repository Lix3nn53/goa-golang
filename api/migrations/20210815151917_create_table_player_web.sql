-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_web` (
     `uuid` varchar(40) NULL,
     `email` varchar(45) NULL,
     `credits` smallint NOT NULL,
     `sessions` text NULL,
     `twitch_id` varchar(40) NOT NULL,
     `discord_id` varchar(40) NOT NULL,
     `google_id` varchar(40) NOT NULL,
     KEY `fkIdx_188` (`uuid`),
     CONSTRAINT `FK_187` FOREIGN KEY `fkIdx_188` (`uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd
-- +goose Down
DROP TABLE goa_player_web;