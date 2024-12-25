package postgres

import (
	"billy/models"
	"billy/storage"
	"database/sql"
	"errors"
	"strings"
)

func (store *PostgresStorage) CreateUser(user models.User) (int, error) {
	var id int
	query := `INSERT INTO users (name, login, pass, plan, role) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := store.db.QueryRow(query, user.Name, user.Login, user.Pass, user.Plan, user.Role).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "users_login_key") {
			return -1, errors.New("login exists")
		}
		return -1, errors.New("error creating user")
	}
	return id, nil
}

func (store *PostgresStorage) GetUsers() ([]models.User, error) {
	var users []models.User
	query := `SELECT id, name, login, pass, plan, role, active FROM users;`
	rows, err := store.db.Query(query)
	if err != nil {
		return nil, errors.New("error fetching users from the database")
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Login, &user.Pass,
			&user.Plan, &user.Role, &user.Active); err != nil {
			return nil, errors.New("error scanning user data")
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("error iterating over user rows")
	}
	return users, nil
}

func (store *PostgresStorage) GetUserById(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, login, pass, plan, role, active FROM users WHERE id = $1;`
	err := store.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Login,
		&user.Pass, &user.Plan, &user.Role, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (store *PostgresStorage) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, login, pass, plan, role, active FROM users WHERE login = $1;`
	err := store.db.QueryRow(query, login).Scan(&user.ID, &user.Name, &user.Login,
		&user.Pass, &user.Plan, &user.Role, &user.Active)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
