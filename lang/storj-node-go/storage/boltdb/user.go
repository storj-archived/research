package boltdb

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Uuid     string `json:"uuid"`
}

// CreateUser calls bolt database instance to create user
func (bdb *Client) CreateUser(key, value []byte) {
	bdb.UsersBucket.Put(key, value)
}

func (bdb *Client) GetUser(key []byte) {
	bdb.UsersBucket.Get(key)
}

func (bdb *Client) DeleteUser(key []byte) {
	bdb.UsersBucket.Delete(key)
}
