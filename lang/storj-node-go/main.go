package main

import (
	"fmt"
	"github.com/Storj/research/lang/storj-node-go/routes"
	"github.com/Storj/research/lang/storj-node-go/storage/boltdb"
	"github.com/kataras/iris"
)

// type Contact struct {
// 	Id int64 `json:"id"`
// }

func main() {

	bdb, err := boltdb.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer bdb.DB.Close()

	users := routes.Users{DB: bdb}
	app := iris.Default()

	SetRoutes(app, users)

	app.Run(iris.Addr(":8080"))
}

// SetRoutes defines all restful routes on the service
func SetRoutes(app *iris.Application, users routes.Users) {
	app.Post("/users", users.CreateUser)
	app.Delete("/users/:id", db.DeleteUser)
}
