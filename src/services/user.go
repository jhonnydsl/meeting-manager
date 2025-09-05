package services

import (
	"errors"
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

// CreateUser handles user creation with basic validation and password hashing.
func (service *UserService) CreateUser(user dtos.UserInput) (dtos.UserOutput, error) {
	// Basic validation: email and name cannot be empty.
	if user.Email == "" || user.Name == "" {
		return dtos.UserOutput{}, errors.New("nome e email não podem estar vazios")
	}

	// Check minimum password length (at least 6 characters).
	if utf8.RuneCountInString(user.Password) < 6 {
		return dtos.UserOutput{}, errors.New("senha deve conter no minimo 6 caracteres")
	}

	// Hash the password using bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return dtos.UserOutput{}, err
	}
	hashedPassword := string(hash)
	user.Password = hashedPassword

	// Save the user into the repository.
	return service.UserRepo.CreateUser(user)
}

// GetAllUsers retrieves all users from the repository.
func (service *UserService) GetAllUsers() ([]dtos.UserOutput, error) {
	return service.UserRepo.GetAllUsers()
}

// LoginUser authenticates user credentials and generates a JWT token.
func (service *UserService) LoginUser(login dtos.UserLoginInput) (string, error) {
	userLogin, err := service.UserRepo.GetUserByEmail(login.Email)
	if err != nil {
		return "", fmt.Errorf("email ou senha invalidos")
	}
	
	// Compare provided password with hashed password in DB.
	err = bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(login.Password))
	if err != nil {
		return "", fmt.Errorf("email ou senha invalidos")
	}

	// Create JWT claims with user data and expiration time.
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

// GetUserByID retrieves a single user by its ID
func (service *UserService) GetUserByID(id int) (dtos.UserOutput, error) {
	return service.UserRepo.GetUserByID(id)
}

// DeleteUser removes a user by its ID
func (service *UserService) DeleteUser(id int) error {
	return service.UserRepo.DeleteUser(id)
}