package routes

import (
	"github.com/Storj/research/lang/storj-node-go/db"
	"github.com/boltdb/bolt"
	"github.com/kataras/iris"
	"log"
	"time"
)

func StartDB() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})

	db = users.DB{
		Bucket: b,
	}
}

func SetRoutes(app *iris.Framework) {
	app.Post("/users", db.CreateUser)
	// app.Get("/users", db.ListUsers)
	// app.Delete("/users/:id", db.DeleteUser)
}
