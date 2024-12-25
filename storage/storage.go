package storage

import (
	"billy/models"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Storage interface {
	UserStorage
	UsageStorage
}
type UserStorage interface {
	CreateUser(models.User) (int, error)
	GetUsers() ([]models.User, error)
	GetUserById(int) (*models.User, error)
	GetUserByLogin(string) (*models.User, error)
}
type UsageStorage interface {
	CreateUsage(models.Usage) (int, error)
	GetUsageByUserID(int) (*models.Usage, error)
	EmulateUsageGathering(id int) (int, error)
}
type InvoiceStorage interface{}
