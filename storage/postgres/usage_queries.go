package postgres

import (
	"billy/models"
	"billy/storage"
	"billy/utils"
	"database/sql"
	"errors"
	"log"
)

func (store *PostgresStorage) CreateUsage(usage models.Usage) (int, error) {
	var id int
	query := `INSERT INTO usage (user_id, youtube, netflix, spotify, basic) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := store.db.QueryRow(query, usage.UserID, usage.Youtube, usage.Netflix, usage.Spotify, usage.Basic).Scan(&id)
	if err != nil {
		log.Print(err)
		return -1, errors.New("error creating usage")
	}
	return id, nil
}

func (store *PostgresStorage) GetUsageByUserID(user_id int) (*models.Usage, error) {
	var usage models.Usage
	query := `SELECT id, youtube, netflix, spotify, basic, verified_at FROM usage WHERE user_id = $1 ORDER BY verified_at DESC LIMIT 1;`
	err := store.db.QueryRow(query, user_id).Scan(&usage.ID, &usage.Youtube, &usage.Netflix,
		&usage.Spotify, &usage.Basic, &usage.VerifiedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}
	return &usage, nil
}

// emulates process of gathering data from data-house
func (store *PostgresStorage) EmulateUsageGathering(id int) (int, error) {
	return store.CreateUsage(utils.GenerateUsage(id))
}
