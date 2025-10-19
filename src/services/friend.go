package services

import (
	"fmt"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
)

type FriendService struct {
	FriendRepo *repository.FriendRepository
}

func (service *FriendService) AddFriend(userID, friendID int) error {
	if userID == friendID {
		return fmt.Errorf("cannot add yourself as a friend")
	}

	return service.FriendRepo.AddFriend(userID, friendID)
}