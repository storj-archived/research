package boltdb

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
func (bdb *Client) DeleteUser(key []byte) {
	bdb.UsersBucket.Delete(key)
}
