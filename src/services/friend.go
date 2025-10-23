package services

import (
	"fmt"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
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

func (service *FriendService) GetFriends(userID int) ([]dtos.FriendOutput, error) {
	return service.FriendRepo.GetFriends(userID)
}

func (service *FriendService) GetFriendsPending(userID int) ([]dtos.FriendOutput, error) {
	return service.FriendRepo.GetFriendsPending(userID)
}

func (service *FriendService) AcceptedFriend(friendID, userID int) error {
	return service.FriendRepo.AcceptFriend(friendID, userID)
}