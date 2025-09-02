package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func (service *UserService) CreateUser(user dtos.UserInput) (dtos.UserOutput, error) {
	if user.Email == "" || user.Name == "" {
		return dtos.UserOutput{}, errors.New("nome e email não podem estar vazios")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.UserOutput{}, err
	}
	hashedPassword := string(hash)
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

	// Cria um novo token JWT.
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