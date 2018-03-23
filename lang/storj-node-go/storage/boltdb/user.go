package boltdb

import (
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

// CreateUser calls bolt database instance to create user
func (bdb *Client) CreateUser(key, value []byte) error {
	return bdb.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))

		return b.Put(key, value)
	})
}

func (bdb *Client) GetUser(key []byte) {
}

func (bdb *Client) DeleteUser(key []byte) {
}
