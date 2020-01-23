package cache

import (
	"github.com/shadkain/db_hw/internal/vars"
	"strings"
	"sync"
)

type UserCache struct {
	idToNickname   map[int]string
	idToNicknameMu sync.RWMutex

	nicknameToID   map[string]int
	nicknameToIDMu sync.RWMutex
}

func NewUserCache() *UserCache {
	return &UserCache{
		idToNickname: make(map[int]string),
		nicknameToID: make(map[string]int),
	}
}

func (this *UserCache) GetIDByNick(nickname string) (int, error) {
	this.nicknameToIDMu.RLock()
	id, ok := this.nicknameToID[strings.ToLower(nickname)]
	this.nicknameToIDMu.RUnlock()

	if !ok {
		return 0, vars.ErrNotFound
	}

	return id, nil
}

func (this *UserCache) GetNicknameByID(id int) (string, error) {
	this.idToNicknameMu.RLock()
	nickname, ok := this.idToNickname[id]
	this.idToNicknameMu.RUnlock()

	if !ok {
		return "", vars.ErrNotFound
	}

	return nickname, nil
}

func (this *UserCache) GetNicknameCaseInsensitive(nickname string) (string, error) {
	this.nicknameToIDMu.RLock()
	id, ok := this.nicknameToID[strings.ToLower(nickname)]
	this.nicknameToIDMu.RUnlock()

	if !ok {
		return "", vars.ErrNotFound
	}

	return this.GetNicknameByID(id)
}

func (this *UserCache) Set(id int, nickname string) {
	this.setIDByNickname(nickname, id)
	this.setNicknameByID(id, nickname)
}

func (this *UserCache) setIDByNickname(nickname string, id int) {
	this.nicknameToIDMu.Lock()
	this.nicknameToID[strings.ToLower(nickname)] = id
	this.nicknameToIDMu.Unlock()
}

func (this *UserCache) setNicknameByID(id int, nickname string) {
	this.idToNicknameMu.Lock()
	this.idToNickname[id] = nickname
	this.idToNicknameMu.Unlock()
}
