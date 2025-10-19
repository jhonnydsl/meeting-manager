package repository

type TableRepository struct{}

func (r *TableRepository) CreateTableUsers() error {
	query := `CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) CreateTableReunioes() error {
	query := `CREATE TABLE IF NOT EXISTS reunioes (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description TEXT,
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	owner_id INT REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) CreateTableConvites() error {
	query := `CREATE TABLE IF NOT EXISTS convites (
	id SERIAL PRIMARY KEY,
	reuniao_id INT REFERENCES reunioes(id),
	user_id INT REFERENCES users(id),
	status VARCHAR(20) DEFAULT 'pending',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (r *TableRepository) CreateTableFriends() error {
	query := `CREATE TABLE IF NOT EXISTS friends (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	friend_id INT REFERENCES users(id),
	status VARCHAR(20) DEFAULT 'pending'
	)`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}