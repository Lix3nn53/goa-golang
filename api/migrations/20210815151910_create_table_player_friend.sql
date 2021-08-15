-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_friend` (
     `uuid`        varchar(40) NOT NULL ,
     `friend_uuid` varchar(40) NOT NULL ,

     UNIQUE KEY `Ind_89` (`friend_uuid`, `uuid`),
     KEY `fkIdx_22` (`uuid`),
     CONSTRAINT `FK_22` FOREIGN KEY `fkIdx_22` (`uuid`) REFERENCES `goa_player` (`uuid`),
     KEY `fkIdx_25` (`friend_uuid`),
     CONSTRAINT `FK_25` FOREIGN KEY `fkIdx_25` (`friend_uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_player_friend;