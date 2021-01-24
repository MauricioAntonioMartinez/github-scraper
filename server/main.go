package main

import (
	"log"
	"time"

	"github.com/MauricioAntonioMartinez/github-scraper/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/helmet/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/storage/sqlite3"
	"github.com/qinains/fastergoding"
)

func main() {
	fastergoding.Run()
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(400).JSON(fiber.Map{
				"success": false,
				"message": err.Error()},
			)
		},
	})

	app.Use(helmet.New())

	jwtCheck := jwtware.New(jwtware.Config{
		SigningKey: []byte("thisismysecret"),
	})

	baseLimiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 30 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"message": "You are rate limited.",
			})
		},
		Storage: sqlite3.New(),
	})

	withAuth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"mcuve": "mcuve",
		},
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			if user == "mcuve" && pass == "mcuve" {
				return true
			}
			return true
		},
		ContextUsername: "_user",
		ContextPassword: "_password",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${locals:mcuve-key} ${status} - ${method} ${path}\n",
		TimeFormat: "18-Jan-2021",
	}))

	app.Use(requestid.New(requestid.Config{
		Header:     "X-CUSTOM-HEADER",
		ContextKey: "mcuve-key",
		Generator: func() string {
			return utils.UUID()
		},
	}))

	app.Use(recover.New())

	routes.AuthRoutes(app, jwtCheck, baseLimiter)
	routes.ChatRoutes(app)
	routes.BasicsRoutes(app, withAuth, jwtCheck, baseLimiter)

	app.Static("/files", "./public")

	log.Fatal(app.Listen(":4000"))
}

