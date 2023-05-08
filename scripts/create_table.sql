create schema codechallenge;

drop table if exists "codechallenge"."users";
create table "codechallenge"."users" (
    "id" bigserial primary key not null,
    "username" varchar(50) not null,
    "password" varchar(255) not null,
    "created_at" timestamptz null,
    "updated_at" timestamptz null,
    "deleted_at" timestamptz null
);


drop table if exists "codechallenge"."metadata";
create table "codechallenge"."metadata" (
    "id" bigserial primary key not null,
    "name" varchar(100) not null,
    "type" varchar(50) not null,
    "size" bigserial not null,
    "location" varchar(255) not null,
    "created_at" timestamptz null,
    "updated_at" timestamptz null,
    "deleted_at" timestamptz null
);
