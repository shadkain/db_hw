package usecase

import (
	"github.com/shadkain/db_hw/internal/reqmodels"
)

func (this *usecaseImpl) GetStatus() (s reqmodels.Status, err error) {
	forum, err := this.repo.CountForums()
	if err != nil {
		return
	}

	post, err := this.repo.CountPosts()
	if err != nil {
		return
	}

	thread, err := this.repo.CountThreads()
	if err != nil {
		return
	}

	user, err := this.repo.CountUsers()
	if err != nil {
		return
	}

	s = reqmodels.Status{
		Forum:  forum,
		Post:   post,
		Thread: thread,
		User:   user,
	}

	return
}

func (this *usecaseImpl) Clear() error {
	return this.repo.Clear()
}
