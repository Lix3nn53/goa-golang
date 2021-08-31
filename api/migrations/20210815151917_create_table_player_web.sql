-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_web` (
     `uuid`        varchar(40) NULL ,
     `google_id`      varchar(40) NULL UNIQUE,
     `twitch_id`      varchar(40) NULL UNIQUE,
     `discord_id`      varchar(40) NULL UNIQUE,
     `email`       varchar(45) NULL UNIQUE,
     `credits`     smallint DEFAULT 0,
     `sessions`    text NULL ,

     KEY `fkIdx_188` (`uuid`),
     CONSTRAINT `FK_187` FOREIGN KEY `fkIdx_188` (`uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_player_web;
