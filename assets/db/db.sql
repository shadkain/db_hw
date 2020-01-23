CREATE EXTENSION IF NOT EXISTS citext;

-- User

CREATE TABLE "user" (
    "id" serial,
    "nickname" citext not null primary key,
    "email" citext not null unique,
    "fullname" text not null,
    "about" text not null default ''
);

-- Forum

CREATE TABLE "forum" (
    "id" serial,
    "slug" citext not null primary key,
    "title" text not null,
    "posts" int not null default 0,
    "threads" int not null default 0,
    "user" varchar not null
);

-- Thread

CREATE TABLE "thread" (
    "id" serial primary key,
    "slug" citext not null,
    "title" text not null,
    "message" text not null,
    "votes" int not null default 0,
    "created" timestamptz not null,
    "forum" varchar not null,
    "author" varchar not null
);
CREATE INDEX ON "thread" ("slug");
CREATE INDEX ON "thread" ("created", "forum");
CREATE INDEX ON "thread" ("forum", "author");

-- Post

CREATE TABLE "post" (
    "id" serial primary key,
    "parent" int not null,
    "path" text not null default '',
    "thread" int not null,
    "message" text not null,
    "isEdited" bool not null default false,
    "created" timestamptz not null,
    "forum" varchar not null,
    "author" varchar not null
);
CREATE INDEX ON "post" ("thread");
CREATE INDEX ON "post" (substring("path",1,7));
CREATE INDEX ON "post" ("forum", "author");

-- Vote

CREATE TABLE "vote" (
    "id" serial primary key,
    "thread" int not null,
    "nickname" varchar not null,
    "voice" int not null
);
CREATE INDEX ON "vote" ("thread", "nickname");

-- Forum-user

CREATE TABLE "forum_user" (
    "forum" varchar not null,
    "user" varchar not null
);
CREATE UNIQUE INDEX ON "forum_user" ("user", "forum");

-- Functions & triggers

CREATE FUNCTION inc_forum_thread() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE forum SET threads = threads + 1 WHERE slug=NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION add_forum_user() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO forum_user (forum, "user") VALUES (NEW.forum, NEW.author) ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE inc_forum_thread();

CREATE TRIGGER forum_user
    AFTER INSERT
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();

CREATE TRIGGER forum_user
    AFTER INSERT
    ON thread
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();
