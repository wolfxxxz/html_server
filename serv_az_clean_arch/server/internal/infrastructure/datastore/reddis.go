package datastore

import "server/internal/domain/models"

type HashDB struct {
	DB map[string]*models.User
}

func InitHashDB() *HashDB {

	var hashTableUsers = make(map[string]*models.User)
	return &HashDB{DB: hashTableUsers}
}
