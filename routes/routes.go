package routes

import (
	"database1/config"
	"database1/controller"
	"database1/middleware"
	"log"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	//socket route
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(message))
	}))

	config.GoogleConfig()
	config.GithubConfig()

	user := app.Group("/user")
	user.Post("/register", controller.Register)
	user.Post("/login", controller.Login)
	user.Get("/google_login", controller.GoogleLogin)
	user.Get("/google_callback", controller.GoogleCallback)
	user.Get("/github_login", controller.GithubLogin)
	user.Get("/github_callback", controller.GithubCallback)

	category := app.Group("/category")
	category.Get("/", middleware.CheckAuth, controller.GetCategory)
	category.Post("/", middleware.CheckAuth, controller.AddCategory)
	category.Put("/:id", middleware.CheckAuth, controller.UpdateCategory)
	category.Delete("/:id", middleware.CheckAuth, controller.DeleteCategory)

	product := app.Group("/product")
	product.Get("/", middleware.CheckAuth, controller.GetProduct)
	product.Post("/", middleware.CheckAuth, controller.AddProduct)
	product.Put("/:id", middleware.CheckAuth, controller.UpdateProduct)
	product.Delete("/:id", middleware.CheckAuth, controller.DeleteProduct)

	cart := app.Group("/cart")
	cart.Post("/", middleware.CheckAuth, controller.AddProductToCart)

}
