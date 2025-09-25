package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/utils"
	"golang.org/x/crypto/bcrypt"
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
		return "", fmt.Errorf("email ou senha invalidos")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(login.Password))
	if err != nil {
		return "", fmt.Errorf("email ou senha invalidos")
	}

	claims := jwt.MapClaims{
		"user_id": userLogin.ID,
		"email": userLogin.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create and sign a new JWT token using secret key.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token: %w", err)
	}

	return tokenStr, nil
}

func (service *UserService) GetUserByID(id int) (dtos.UserOutput, error) {
	return service.UserRepo.GetUserByID(id)
}

func (service *UserService) DeleteUser(id int) error {
	return service.UserRepo.DeleteUser(id)
}