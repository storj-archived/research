package routes

import (
	"github.com/Storj/research/lang/storj-node-go/storage/boltdb"
	"github.com/google/uuid"
	"github.com/kataras/iris"
	"log"
)

// Users contains items needed to process requests to the user namespace
type Users struct {
	DB *boltdb.Client
}

// CreateUser instantiates a new user
func (u *Users) CreateUser(ctx iris.Context) {
	user := boltdb.User{
		Id:       uuid.New(),
		Username: ctx.Params().Get("id"),
		Email:    `dece@trali.zzd`,
	}

	if err := ctx.ReadJSON(user); err != nil {
		ctx.JSON(iris.StatusNotAcceptable)
	}

	u.DB.CreateUser(user)
}

func (u *Users) GetUser(ctx iris.Context) {
	userId := ctx.Params().Get("id")
	userInfo, err := u.DB.GetUser([]byte(userId))
	if err != nil {
		log.Println(err)
	}

	ctx.Writef("%s's info is: %s", userId, userInfo)
}

func (u *Users) EditUser(ctx iris.Context) {
}

// DeleteUser deletes a user key/value from users bucket
func (u *Users) DeleteUser(ctx iris.Context) {
	user := &boltdb.User{}
	u.DB.DeleteUser([]byte(user.Username))
}
