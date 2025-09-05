package repository

import (
	"github.com/jhonnydsl/gerenciamento-de-reunioes/src/dtos"
)

type UserRepository struct{}

// CreateUser inserts a new user into the database and returns the created user.
func (r *UserRepository) CreateUser(user dtos.UserInput) (dtos.UserOutput, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email`
	var createdUser dtos.UserOutput

	// Execute the query and scan the returned values into createdUser.
	err := DB.QueryRow(query, user.Name, user.Email, user.Password).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email)
	if err != nil {
		return dtos.UserOutput{}, err
	}

	return createdUser, nil
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers() ([]dtos.UserOutput, error) {
	query := `SELECT id, name, email FROM users`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []dtos.UserOutput

	// Iterate through rows and append each user to the list.
	for rows.Next() {
		var u dtos.UserOutput

		err = rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}

		lista = append(lista, u)
	}

	return lista, nil
}

// GetUserByEmail retrieves a user (including password hash) by email.
func (r *UserRepository) GetUserByEmail(email string) (dtos.UserLogin, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	var login dtos.UserLogin

	row := DB.QueryRow(query, email)

	// Scan result into user struct.
	if err := row.Scan(&login.ID, &login.Name, &login.Email, &login.Password); err != nil {
		return dtos.UserLogin{}, err
	}

	return login, nil
}

// GetUserByID retrieves a single user (without password) by ID.
func (r *UserRepository) GetUserByID(id int) (dtos.UserOutput, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	var user dtos.UserOutput

	row := DB.QueryRow(query, id)

	// Scan result into user struct.
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return dtos.UserOutput{}, err
	}

	return user, nil
}

// DeleteUser removes a user by ID.
func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
