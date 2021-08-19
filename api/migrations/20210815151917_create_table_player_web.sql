-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player_web` (
     `uuid`        varchar(40) NOT NULL ,
     `email`       varchar(45) NULL ,
     `mc_username` varchar(20) NOT NULL ,
     `credits`     smallint NOT NULL ,
     `sessions`    text NULL ,

     PRIMARY KEY (`uuid`),
     KEY `fkIdx_38_clone` (`uuid`),
     CONSTRAINT `FK_38_clone` FOREIGN KEY `fkIdx_38_clone` (`uuid`) REFERENCES `goa_player` (`uuid`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_player_web;