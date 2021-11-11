-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_player` (
     `uuid`              varchar(40) NOT NULL ,
     `daily_last_date`   date NULL ,
     `staff_rank`        varchar(20) NULL ,
     `premium_rank`      varchar(20) NULL ,
     `premium_rank_date` date NULL ,
     `storage_personal`  mediumtext NULL ,
     `storage_bazaar`    mediumtext NULL ,
     `storage_premium`   mediumtext NULL ,
     `lang`              varchar(20) NULL ,
     
     PRIMARY KEY (`uuid`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_player;