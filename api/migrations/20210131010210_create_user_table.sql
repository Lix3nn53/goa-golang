-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_guild" (
     "name"         varchar(20) NOT NULL ,
     "tag"          varchar(5) NOT NULL ,
     "war_point"    smallint NULL ,
     "announcement" tinytext NULL ,
     "hall_level"   smallint NOT NULL ,
     "bank_level"   smallint NOT NULL ,
     "lab_level"    smallint NOT NULL ,
     "storage"      text NULL ,

     PRIMARY KEY ("name")
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_player" (
     "uuid"              varchar(40) NOT NULL ,
     "daily_last_date"   date NULL ,
     "staff_rank"        varchar(20) NULL ,
     "premium_rank"      varchar(20) NULL ,
     "premium_rank_date" date NULL ,
     "storage_personal"  mediumtext NULL ,
     "storage_bazaar"    mediumtext NULL ,
     "storage_premium"   mediumtext NULL ,
     
     PRIMARY KEY ("uuid")
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_player_character" (
     "character_no"         smallint NOT NULL ,
     "uuid"                 varchar(40) NOT NULL ,
     "off_hand"             text NULL ,
     "slot_parrot"          text NULL ,
     "slot_necklace"        text NULL ,
     "slot_ring"            text NULL ,
     "slot_earring"         text NULL ,
     "slot_glove"           text NULL ,
     "slot_pet"             text NULL ,
     "chat_tag"             varchar(45) NULL ,
     "crafting_experiences" text NOT NULL ,
     "inventory"            mediumtext NOT NULL ,
     "turnedinquests"       text NULL ,
     "activequests"         text NULL ,
     "location"             text NOT NULL ,
     "armor_content"        text NOT NULL ,
     "rpg_class"            varchar(45) NOT NULL ,
     "unlocked_classes"     mediumtext NULL ,
     "totalexp"             int NOT NULL ,
     "attr_one"             smallint NOT NULL ,
     "attr_two"             smallint NOT NULL ,
     "attr_three"           smallint NOT NULL ,
     "attr_four"            smallint NOT NULL ,
     "attr_five"            smallint NOT NULL ,
     "skill_one"            smallint NOT NULL ,
     "skill_two"            smallint NOT NULL ,
     "skill_three"          smallint NOT NULL ,
     "skill_passive"        smallint NOT NULL ,
     "skill_ultimate"       smallint NOT NULL ,

     UNIQUE KEY "Ind_88" ("uuid", "character_no"),
     KEY "fkIdx_55" ("uuid"),
     CONSTRAINT "FK_55" FOREIGN KEY "fkIdx_55" ("uuid") REFERENCES "goa_player" ("uuid")
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_player_friend" (
     "uuid"        varchar(40) NOT NULL ,
     "friend_uuid" varchar(40) NOT NULL ,

     UNIQUE KEY "Ind_89" ("friend_uuid", "uuid"),
     KEY "fkIdx_22" ("uuid"),
     CONSTRAINT "FK_22" FOREIGN KEY "fkIdx_22" ("uuid") REFERENCES "goa_player" ("uuid"),
     KEY "fkIdx_25" ("friend_uuid"),
     CONSTRAINT "FK_25" FOREIGN KEY "fkIdx_25" ("friend_uuid") REFERENCES "goa_player" ("uuid")
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_player_guild" (
     "uuid" varchar(40) NOT NULL ,
     "name" varchar(20) NOT NULL ,
     "rank" varchar(20) NOT NULL ,
     
     PRIMARY KEY ("uuid"),
     KEY "fkIdx_38" ("uuid"),
     CONSTRAINT "FK_38" FOREIGN KEY "fkIdx_38" ("uuid") REFERENCES "goa_player" ("uuid"),
     KEY "fkIdx_41" ("name"),
     CONSTRAINT "FK_41" FOREIGN KEY "fkIdx_41" ("name") REFERENCES "goa_guild" ("name")
);
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "goa_web_player" (
     "uuid"        varchar(40) NOT NULL ,
     "email"       varchar(45) NULL ,
     "mc_username" varchar(20) NOT NULL ,
     "credits"     smallint NOT NULL ,

     PRIMARY KEY ("uuid"),
     KEY "fkIdx_38_clone" ("uuid"),
     CONSTRAINT "FK_38_clone" FOREIGN KEY "fkIdx_38_clone" ("uuid") REFERENCES "goa_player" ("uuid")
);
-- +goose StatementEnd