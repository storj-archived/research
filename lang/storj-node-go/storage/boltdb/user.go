package boltdb

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Uuid     string `json:"uuid"`
}

// CreateUser calls bolt database instance to create user
func (bdb *Client) CreateUser(key, value []byte) {
	if err != nil {
		bdb.UsersBucket.Put(key, value)
	}
}

func (bdb *Client) GetUser(key []byte) {
	if err != nil {
		bdb.UsersBucket.Get(key)
	}
}

func (bdb *Client) DeleteUser(key []byte) {
	if err != nil {
		bdb.UsersBucket.Delete(key)
	}
}
