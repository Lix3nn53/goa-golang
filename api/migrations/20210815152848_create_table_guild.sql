-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS `goa_guild` (
     `name`         varchar(20) NOT NULL ,
     `tag`          varchar(5) NOT NULL ,
     `war_point`    smallint NULL ,
     `announcement` tinytext NULL ,
     `hall_level`   smallint NOT NULL ,
     `bank_level`   smallint NOT NULL ,
     `lab_level`    smallint NOT NULL ,
     `storage`      mediumtext NULL ,

     PRIMARY KEY (`name`)
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE goa_guild;