package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/qinains/fastergoding"
)

type UserResponse struct {
	name string
	age int32
	degree string
}


func main() {
	fastergoding.Run()
	app := fiber.New(fiber.Config{
		ETag: true,
		ErrorHandler: func(ctx *fiber.Ctx,err error)error{

			
			return ctx.Status(400).JSON(fiber.Map{
				"success":false,
				"message":err.Error(),
				"errors":fiber.Map{"type":"unknown","field":"This ass"}},
			)
		},
		
	})



	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",

	}))


	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:mcuve-key} ${status} - ${method} ${path}\n",
		TimeFormat: "18-Jan-2021",

	}),)

	app.Use(requestid.New(requestid.Config{
		Header: "X-CUSTOM-HEADER",
		ContextKey:"mcuve-key",
		Generator: func()string{
			return utils.UUID()
		},
	}))

	app.Use(recover.New())

	

	withAuth :=basicauth.New(basicauth.Config{
		Users : map[string]string{
			"mcuve":"mcuve",
		},
		Realm:"Forbidden",
		Authorizer: func(user, pass string) bool{
		
			if user == "mcuve" && pass == "mcuve" {
				return true
			} 
			return true
		},
		// Unauthorized: func(c *fiber.Ctx) error  { 
		// 	fmt.Println("something")
		// 	return c.SendFile("./images/docker.png")
		// },
		ContextUsername: "_user",
		ContextPassword: "_password",

	})


	app.Static("/files","./public")

    app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println( c.Get("X-CUSTOM-HEADER"))
        return c.SendString("Hello, World 👋!")
	})

	app.Get("/error",func(c *fiber.Ctx)error{
		panic("Something went wrong")
		
	})

	
	app.Get("/:file.:ext",func(c *fiber.Ctx) error  { 

			msg := fmt.Sprintf("This is a file with an extension: %s",c.Params("ext"))
			return c.SendString(msg)
	})
	
	
	app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
        return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	app.Use("/api",withAuth, func(c *fiber.Ctx) error {
		return c.Next()
		// return c.JSON(UserResponse{name:"something",age: 21,degree: "engineer"})
	})

	app.Get("api/list",func(c *fiber.Ctx) error { 
		
		c.SendStatus(201)
		app := c.App()
		app.Static("/dynamic","/images")

		return c.SendString("happy")
	})

	app.Use("/next/1",func(c *fiber.Ctx)error { 

		fmt.Println("Thi is a handler")
		return c.JSON(`{ "name":"simon" } `)
	})
	
	app.Get("/:name/:age?/:gender?",func (c *fiber.Ctx) error {
		// the ? means this param is optional
		msg := fmt.Sprintf(" %s is %s years old",c.Params("name"),c.Params("age"))
		return c.SendString(msg)
	})



   log.Fatal( app.Listen(":3000"))
}

func willError(rut chan string){
	rut <- "Hello"
}



