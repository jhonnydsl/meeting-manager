package services

import (
	"errors"
	"fmt"

	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func (service *UserService) CreateUser(user dtos.UserInput) (dtos.UserOutput, error) {
	if err := utils.ValidateUserInput(user); err != nil {
		return dtos.UserOutput{}, err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return dtos.UserOutput{}, err
	}
	user.Password = hashedPassword

	return service.UserRepo.CreateUser(user)
}

func (service *UserService) GetAllUsers() ([]dtos.UserOutput, error) {
	return service.UserRepo.GetAllUsers()
}

func (service *UserService) LoginUser(login dtos.UserLoginInput) (string, error) {
	userLogin, err := service.UserRepo.GetUserByEmail(login.Email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}
	
	if err := utils.CheckPassword(userLogin.Password, login.Password); err != nil {
		return "", errors.New("invalid email or password")
	}

	tokenStr, err := utils.GenerateJWT(userLogin.ID, userLogin.Email)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenStr, nil
}

func (service *UserService) GetUserByID(id int) (dtos.UserOutput, error) {
	return service.UserRepo.GetUserByID(id)
}

func (service *UserService) DeleteUser(id int) error {
	return service.UserRepo.DeleteUser(id)
}