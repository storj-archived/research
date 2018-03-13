package boltdb

import (
	"fmt"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"uuid"`
	Uuid     string `json:"uuid"`
}

// CreateUser calls bolt database instance to create user
func (bdb *Client) CreateUser(key, value []byte) {
	bdb.UsersBucket.Put(key, value)
}

// DeleteUser calls bolt database instance to delete user
func (b *Bucket) DeleteUser(key []byte) {
	b.Delete(key)
}
