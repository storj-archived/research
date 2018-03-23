package routes

import (
	"encoding/json"
	"fmt"
	"github.com/Storj/research/lang/storj-node-go/storage/boltdb"
	"github.com/google/uuid"
	"github.com/kataras/iris"
)

// Users contains items needed to process requests to the user namespace
type Users struct {
	DB *boltdb.Client
}

// CreateUser instantiates a new user
func (u *Users) CreateUser(ctx iris.Context) {
	user := &boltdb.User{
		Id:       uuid.New(),
		Username: ctx.Params().Get("id"),
		Email:    `dece@trali.zzd`,
	}

	if err := ctx.ReadJSON(user); err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	usernameKey := []byte(user.Username)

	u.DB.CreateUser(usernameKey, userBytes)
}

func (u *Users) GetUser(ctx iris.Context) {
	user := &boltdb.User{}
	u.DB.GetUser([]byte(user.Username))
}

// DeleteUser deletes a user key/value from users bucket
func (u *Users) DeleteUser(ctx iris.Context) {
	user := &boltdb.User{}
	u.DB.DeleteUser([]byte(user.Username))
}
