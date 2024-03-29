package routes
import (
	"github.com/gofiber/fiber/v2"
	"github.com/MichaelYoung87/backend-go-project-remake/controllers"	
)
func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
	app.Post("/api/checkIfUsernameExists", controllers.CheckIfUsernameExists)
	app.Post("/api/uploadProfilePicture", controllers.UploadProfilePicture)
}