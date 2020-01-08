-- Users

CREATE TABLE users (
    id serial not null primary key,
    nickname varchar not null,
    email varchar not null,
    fullname varchar,
    about text
);

CREATE UNIQUE INDEX user_id_uindex ON users USING btree (id);
CREATE UNIQUE INDEX user_nickname_uindex ON users USING btree (nickname);
CREATE UNIQUE INDEX user_email_uindex ON users USING btree (email);

-- Forums

CREATE TABLE forums (
    id serial not null primary key,
    slug varchar not null,
    title varchar not null,
    posts int default 0 not null,
    threads int default 0 not null,
    creator varchar not null references users(nickname)
);

CREATE UNIQUE INDEX forum_id_uindex ON forums USING btree (id);
CREATE UNIQUE INDEX forum_slug_uindex ON forums USING btree (slug);

-- Threads

CREATE TABLE threads (
    id serial not null primary key,
    slug varchar,
    title varchar not null,
    message text not null,
    created timestamptz not null default now(),
    votes int default 0 not null,
    forum varchar not null references forums(slug),
    author varchar not null references users(nickname)
);

CREATE UNIQUE INDEX table_name_id_uindex ON threads USING btree (id);
CREATE UNIQUE INDEX table_name_slug_uindex ON threads USING btree (slug);

-- Posts

CREATE TABLE posts (
    id serial not null primary key,
    created timestamptz not null default '1970-01-01 00:00:00+00'::timestamptz,
    isedited boolean default false not null,
    message text not null,
    parent int default 0 not null,
    forum varchar not null references forums(slug),
    thread int not null references threads(id),
    author varchar not null references users(nickname)
);

CREATE UNIQUE INDEX post_id_uindex ON posts USING btree (id);

-- Votes

CREATE TABLE votes (
    id serial not null primary key,
    voice smallint not null,
    thread int not null references threads(id),
    nickname varchar not null references users(nickname)
);

ALTER TABLE ONLY votes ADD CONSTRAINT subscribe_subscriber_id_followee_id_key UNIQUE (nickname, thread);
CREATE UNIQUE INDEX vote_id_uindex ON votes USING btree (id);

-- Functions

CREATE OR REPLACE FUNCTION vote_add() RETURNS TRIGGER AS $emp_audit$
    BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NULL;
    END;
    $emp_audit$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION thread_add() RETURNS TRIGGER AS $emp_audit$
    BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
    END;
$emp_audit$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION post_add() RETURNS TRIGGER AS $emp_audit$
    BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NULL;
    END;
$emp_audit$ LANGUAGE plpgsql;

-- Triggers

CREATE TRIGGER vote_insert
    AFTER INSERT
    ON votes
    FOR EACH ROW EXECUTE PROCEDURE vote_add(vote);


CREATE TRIGGER thread_insert
    AFTER INSERT
    ON threads
    FOR EACH ROW EXECUTE PROCEDURE thread_add();


CREATE TRIGGER post_insert
    AFTER INSERT
    ON posts
    FOR EACH ROW EXECUTE PROCEDURE post_add();
