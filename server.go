package main

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/miceremwirigi/go-fiber-jwt/databases"
	"github.com/miceremwirigi/go-fiber-jwt/models"
	"github.com/miceremwirigi/go-fiber-jwt/apis"

)

func main() {
	db, _ := db.InitializeDB()

	//Auto migrate
	db.AutoMigrate(&models.User{})

	app := fiber.New()

	h := apis.Handler{DB: db}
	h.SetupRoutes(app)
	
	if err := app.Listen(":5000"); err != nil {
		panic(err)
	}
}
