package routes

import (
	"encoding/json"
	"github.com/Storj/research/lang/storj-node-go/storage/boltdb"
	"github.com/kataras/iris"
	uuid "github.com/satori/go.uuid"
)

// Users contains items needed to process requests to the user namespace
type Users struct {
	DB *boltdb.Client
}

// CreateUser instantiates a new user
func (u *Users) CreateUser(ctx iris.Context) {
	user := &boltdb.User{}

	if err := ctx.ReadJSON(user); err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	uu, err := uuid.NewV4()
	var uuidBytes = []byte(uu.String())

	u.DB.CreateUser(uuidBytes, userBytes)
}

// DeleteUser deletes a user key/value from users bucket
func (u *Users) DeleteUser(ctx iris.Context) {
	// magically get uuid from id
	// u.DB.DeleteUser([]byte(uuid))
}
