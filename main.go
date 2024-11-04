package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {

	dsn := "root@tcp(127.0.0.1:3306)/node_crud_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mysql ")
	// สร้างแอป Fiber ใหม่

	app := fiber.New()
	store := session.New()

	app.Post("/login", func(c *fiber.Ctx) error {

		username := c.FormValue("username")
		password := c.FormValue("passwordd")

		if username == "admin" && password == "123456" {
			sess, err := store.Get(c)
			if err != nil {
				return err
			}

			sess.Set("authentication", true)
			sess.Set("user", username)

			if err := sess.Save(); err != nil {
				return err
			}
			return c.SendString("loggen in  successfully!")
		}
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid credentials")
	})

	app.Get("/api/user", func(c *fiber.Ctx) error {
		var users []User

		rows, err := db.Query("SELECT id,name,email FROM bdname")
		if err != nil {
			return c.Status(500).SendString("Database query error")
		}
		defer rows.Close()

		for rows.Next() {
			var user User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
				return c.Status(500).SendString("Error scening data")
			}
			users = append(users, user)
		}

		return c.JSON(users)
	})

	// รันแอปบนพอร์ต 3000
	app.Listen(":3000")
}
