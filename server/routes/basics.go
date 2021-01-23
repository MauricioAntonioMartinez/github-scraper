package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//BasicRoutes the basics of the fiber framework
func BasicsRoutes(app *fiber.App, withAuth, jwtCheck, baseLimiter func(*fiber.Ctx) error) {

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(c.Get("X-CUSTOM-HEADER"))
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		panic("Something went wrong")

	})

	app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
		return c.SendString(msg) // this only functions with "-"
	})

	app.Use("/api", withAuth, func(c *fiber.Ctx) error {
		return c.Next()
	})

	app.Get("api/list", withAuth, func(c *fiber.Ctx) error {
		c.SendStatus(201)
		app := c.App()
		// dynamically serves folders
		app.Static("/dynamic", "/images")

		return c.SendString("happy")
	})

	app.Get("/restricted/user", jwtCheck, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Restriced route"})
	})

	app.Use("/next/sends-json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "This is sends json"})
	})

	app.Use("/ratelimit", baseLimiter, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Rate limit end point.",
		})
	})
	app.Get("/:file.:ext", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("This is a file with an extension: %s", c.Params("ext"))
		return c.SendString(msg)
	})

	app.Get("/:name/:age?/:gender?", func(c *fiber.Ctx) error {
		// the ? means this param is optional
		msg := fmt.Sprintf(" %s is %s years old", c.Params("name"), c.Params("age"))
		return c.SendString(msg)
	})

}
