package main

import (
	"fmt"
	"github.com/Storj/research/lang/storj-node-go/db"
	"github.com/Storj/research/lang/storj-node-go/routes"
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

// type Contact struct {
// 	Id int64 `json:"id"`
// }

func main() {
	app := iris.Default()

	SetRoutes(app)

	// app.Get("/", func(ctx iris.Context) {
	// 	user := make(map[string]string)
	// 	user["username"] = "admin"
	// 	ctx.JSON(user)
	// })

	// app.Get("/users/{id:long min(1)}", func(ctx iris.Context) {
	// 	var user User

	// 	id, _ := ctx.Params().GetInt64("id")

	// 	user.Id = id
	// 	user.Username = "admin"

	// 	ctx.JSON(user)
	// })

	// app.Delete("/users/{id:long min(1)}", func(ctx iris.Context) {
	// 	ctx.JSON(user)
	// })

	// app.Get("/contacts", func(ctx iris.Context) {
	// 	users := make([]string, 0)
	// 	users = append(users, "user1")
	// 	ctx.JSON(users)
	// })

	// app.Get("/contacts/{id:long min(1)}", func(ctx iris.Context) {

	// 	details := make(map[string]string)
	// 	details["id"] = ctx.Params().Get("id")
	// 	details["username"] = "admin"

	// 	ctx.JSON(details)
	// })

	app.Run(iris.Addr(":8080"))
}
