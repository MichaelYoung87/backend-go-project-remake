package main

import (
	"os"

	"github.com/MichaelYoung87/backend-go-project-remake/database"
	"github.com/MichaelYoung87/backend-go-project-remake/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)
func main() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "keys/google/NAMN_PÅ_JSON_HÄMTAD_FRÅN_GOOGLE_CLOUD.json")
	database.Connect()
		app := fiber.New()
		app.Use(cors.New(cors.Config {
			AllowCredentials: true,
		}))
		routes.Setup(app)
		app.Listen(":8000")
	}