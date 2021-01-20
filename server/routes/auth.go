package routes

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"

	jwt "github.com/form3tech-oss/jwt-go"
)

type InputLogin struct {
	email string
	password string
}

// Validate user login input
func (a InputLogin) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.email, validation.Required, is.Email,),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.password, validation.Required, validation.Length(5, 50)),
	)
}


// AuthRoutes GET,POST,PUT,PATCH /auth
func AuthRoutes(app *fiber.App){
	router := app.Group("/auth")


	router.Get("/login",login)
	
}


func login(c *fiber.Ctx) error {
	req := c.Request() 

	 inputs := InputLogin{}

	if err:= json.Unmarshal(req.Body(),&inputs); err != nil{
		panic("Invalid Format")
	}

	

	err := inputs.Validate()


	if err != nil {
		return c.JSON(err)
	}
	fmt.Println("Not here")

	// pass := c.FormValue("pass")

	// Throws Unauthorized error
	// if user != "john" || pass != "doe" {
	// 	return c.SendStatus(fiber.StatusUnauthorized)
	// }

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "John Doe"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token":t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}