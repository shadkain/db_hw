package query

const (
	// Forums
	InsertForum        = "INSERT INTO forums (slug, title, creator) values ($1, $2, $3) RETURNING id"
	SelectForumsBySlug = "SELECT posts, slug, threads, title, creator FROM forums WHERE lower(slug) = lower($1)"

	// Users
	InsertUser            = "INSERT INTO users (about, email, fullname, nickname) values ($1, $2, $3, $4) RETURNING id"
	SelectUsersByNickname = "SELECT about, email, fullname, nickname FROM users WHERE lower(nickname) = lower($1)"
	SelectUsersByEmail    = "SELECT about, email, fullname, nickname " +
		"FROM users WHERE lower(email) = lower($1);"
	SelectUsersByNicknameOrEmail = "SELECT about, email, fullname, nickname " +
		"FROM users WHERE lower(email) = lower($1) OR lower(nickname) = lower($2);"
	SelectUsersByForumSlug = "SELECT u.about, u.email, u.fullname, u.nickname " +
		"FROM users as u " +
		"WHERE u.nickname IN ( " +
		"SELECT t.author AS nickname " +
		"FROM threads as t " +
		"WHERE lower(t.forum) = lower($1) " +
		"UNION " +
		"SELECT author AS nickname " +
		"FROM posts as p " +
		"WHERE lower(forum) = lower($1) ) " +
		"ORDER BY lower(u.nickname) " +
		"LIMIT 100"
	SelectUsersByForumSlugDesc = "SELECT u.about, u.email, u.fullname, u.nickname " +
		"FROM users as u " +
		"WHERE u.nickname IN ( " +
		"SELECT t.author AS nickname " +
		"FROM threads as t " +
		"WHERE lower(t.forum) = lower($1) " +
		"UNION " +
		"SELECT author AS nickname " +
		"FROM posts as p " +
		"WHERE lower(forum) = lower($1) ) " +
		"ORDER BY u.nickname DESC " +
		"LIMIT 100"
	UpdateUserByNickname = "UPDATE users SET about = $1, email = $2, fullname = $3 WHERE nickname = $4"

	// Threads
	InsertThread                       = "INSERT INTO threads (author, created, message, title, forum) values ($1, $2, $3, $4, $5) RETURNING id"
	InsertThreadWithoutCreated         = "INSERT INTO threads (author, message, title, forum) values ($1, $2, $3, $4) RETURNING id"
	InsertThreadWithSlugWithoutCreated = "INSERT INTO threads (author, message, title, forum, slug) values ($1, $2, $3, $4, $5) RETURNING id"
	InsertThreadWithSlug               = "INSERT INTO threads (author, created, message, title, forum, slug) values ($1, $2, $3, $4, $5, $6) RETURNING id"
	SelectThreadsByForum               = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE lower(forum) = lower($1) ORDER BY created LIMIT $2"
	SelectThreadsByForumSince = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE lower(forum) = lower($1) AND created >= $3 ORDER BY created LIMIT $2"
	SelectThreadsByForumDesc = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE lower(forum) = lower($1) ORDER BY created DESC LIMIT $2"
	SelectThreadsByForumSinceDesc = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE lower(forum) = lower($1) AND created <= $3 ORDER BY created DESC LIMIT $2"
	SelectThreadsBySlug = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE lower(slug) = lower($1)"
	SelectThreadsByID = "SELECT author, created, forum, id, message, slug, title, votes " +
		"FROM threads WHERE id = $1"
	UpdateThreadByID = "UPDATE threads SET message = $1, title = $2 WHERE id = $3"

	// Posts
	InsertPost = "INSERT INTO posts (author, message, parent, thread, forum) " +
		"VALUES ($1, $2, $3, $4, $5) RETURNING id, thread"
	SelectPostsByID = "SELECT author, created, forum, id, isEdited, message, parent, thread " +
		"FROM posts WHERE id = $1"
	SelectPostsByIDThreadID = "SELECT author, created, forum, id, isEdited, message, parent, thread " +
		"FROM posts WHERE id = $1 AND thread = $2"
	SelectPostsFlat = "SELECT author, created, forum, id, isEdited, message, parent, thread " +
		"FROM posts WHERE thread = $1 AND id > $3 ORDER BY id LIMIT $2"
	SelectPostsFlatDesc = "SELECT author, created, forum, id, isEdited, message, parent, thread " +
		"FROM posts WHERE thread = $1 AND id < $3 ORDER BY id DESC LIMIT $2"
	SelectPostsTree = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 and T1.thread = $1 " +
		"UNION " +
		"SELECT T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread from temp1 ORDER BY root, PATH LIMIT $2"
	SelectPostsTreeSince = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1 " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"UNION " +
		"SELECT T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1 " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread from temp1 WHERE id > $3 ORDER BY PATH LIMIT $2"
	SelectPostsTreeDesc = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"UNION " +
		"SELECT  T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread from temp1 WHERE id < $3 ORDER BY PATH DESC LIMIT $2"
	SelectPostsTreeSinceDesc = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (1000000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"UNION " +
		"SELECT  T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST (temp1.PATH ||'->'|| T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON (temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread from temp1 ORDER BY PATH"
	SelectPostsParentTree = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"UNION " +
		"SELECT T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread from temp1 ORDER BY root, PATH"
	SelectPostsParentTreeDesc = "WITH RECURSIVE temp1 (author, created, forum, id, isEdited, message, parent, thread, PATH, LEVEL, root ) AS ( " +
		"SELECT T1.author, T1.created, T1.forum, T1.id, T1.isEdited, T1.message, T1.parent, T1.thread, CAST (10000 + T1.id AS VARCHAR (50)) as PATH, 1, T1.id as root " +
		"FROM posts as T1 WHERE T1.parent = 0 AND T1.thread = $1" +
		"UNION " +
		"SELECT  T2.author, T2.created, T2.forum, T2.id, T2.isEdited, T2.message, T2.parent, T2.thread, CAST ( temp1.PATH ||'->'|| 10000 + T2.id AS VARCHAR(50)), LEVEL + 1, root " +
		"FROM posts as T2 INNER JOIN temp1 ON( temp1.id = T2.parent) " +
		") " +
		"SELECT author, created, forum, id, isEdited, message, parent, thread  from temp1 ORDER BY root desc, PATH"
	UpdatePostByID = "UPDATE posts SET message = $1, isEdited = $2 WHERE id = $3"

	// Votes
	InsertVote = "INSERT INTO votes (nickname, voice, thread) " +
		"VALUES ($1, $2, $3) RETURNING id"
	UpdateVote = "UPDATE votes SET voice = $1 WHERE nickname = $2 AND thread = $3"

	// Service
	SelectStatus = "SELECT " +
		"(SELECT COALESCE(SUM(posts), 0) FROM forums WHERE posts > 0) AS post, " +
		"(SELECT COALESCE(SUM(threads), 0) FROM forums WHERE threads > 0) AS thread, " +
		"(SELECT COUNT(*) FROM users) AS user, " +
		"(SELECT COUNT(*) FROM forums) AS forum"
	ClearAll = "TRUNCATE votes, posts, threads, forums, users RESTART IDENTITY CASCADE"
)
