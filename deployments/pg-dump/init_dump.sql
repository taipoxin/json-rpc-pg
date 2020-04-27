-- Adminer 4.7.6 PostgreSQL dump
CREATE DATABASE "jsonrpc";

\connect "jsonrpc";

DROP TABLE IF EXISTS "posts";
DROP SEQUENCE IF EXISTS posts_id_seq;
CREATE SEQUENCE posts_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."posts" (
    "id" bigint DEFAULT nextval('posts_id_seq') NOT NULL,
    "title" character varying(50) NOT NULL
) WITH (oids = false);

INSERT INTO "posts" ("id", "title") VALUES
(1,	'test');

-- 2020-04-27 11:47:44.484173+00
