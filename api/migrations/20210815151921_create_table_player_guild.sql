-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_guild` (
     `uuid` varchar(40) NOT NULL ,
     `name` varchar(20) NOT NULL ,
     `rank` varchar(20) NOT NULL ,
     
     PRIMARY KEY (`uuid`),
     KEY `fkIdx_38` (`uuid`),
     CONSTRAINT `FK_38` FOREIGN KEY `fkIdx_38` (`uuid`) REFERENCES `goa_player` (`uuid`),
     KEY `fkIdx_41` (`name`),
     CONSTRAINT `FK_41` FOREIGN KEY `fkIdx_41` (`name`) REFERENCES `goa_guild` (`name`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_player_guild;