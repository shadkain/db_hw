package repository

func (this *repositoryImpl) CountForums() (count int, err error) {
	return this.count("forum")
}

func (this *repositoryImpl) CountThreads() (count int, err error) {
	return this.count("thread")
}

func (this *repositoryImpl) CountPosts() (count int, err error) {
	return this.count("post")
}

func (this *repositoryImpl) CountUsers() (count int, err error) {
	return this.count("user")
}

func (this *repositoryImpl) Clear() (err error) {
	_, err = this.db.Exec(`TRUNCATE thread, post, forum, "user", vote, forum_user`)
	return
}

func (this *repositoryImpl) count(table string) (count int, err error) {
	err = this.db.Get(
		&count,
		`SELECT COUNT(*) FROM "`+table+`"`,
	)
	return
}
