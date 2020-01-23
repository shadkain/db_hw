package usecase

import (
	"github.com/shadkain/db_hw/internal/reqmodels"
	"github.com/shadkain/db_hw/internal/models"
)

func (this *usecaseImpl) VoteForThread(threadSlugOrID string, vote reqmodels.Vote) (*models.Thread, error) {
	thread, err := this.repo.GetThreadBySlugOrID(threadSlugOrID)
	if err != nil {
		return nil, err
	}

	trueNickname, err := this.repo.GetUserNickname(vote.Nickname)
	if err != nil {
		return nil, err
	}

	newVotes, err := this.repo.AddThreadVote(thread, trueNickname, vote.Voice)
	thread.Votes = newVotes

	return thread, err
}
