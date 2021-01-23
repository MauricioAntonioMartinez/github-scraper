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
	Email    string
	Password string
}

func (a InputLogin) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required, validation.Length(5, 50)),
	)
}

// AuthRoutes GET,POST,PUT,PATCH /auth
func AuthRoutes(app *fiber.App, jwtCheck, baseLimiter func(*fiber.Ctx) error) {
	router := app.Group("/auth")

	router.Post("/login", login)
	router.Get("/me", baseLimiter, jwtCheck, currentUser, restricted)

}

func login(c *fiber.Ctx) error {
	req := c.Request()

	inputs := InputLogin{}

	if err := json.Unmarshal(req.Body(), &inputs); err != nil {
		panic("Invalid Format")
	}

	err := inputs.Validate()

	if err != nil {
		return c.JSON(err)
	}

	if false {
		// this returns an error
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Mauricio Martinez"
	claims["admin"] = true
	claims["age"] = 21
	claims["expiresIn"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("thisismysecret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t, "expiresIn": 3600 * 72})
}

type data struct{}

func (data) Write(p []byte) (n int, err error) {

	return 2, nil
}

type UserData struct {
	User interface{}
}

func currentUser(c *fiber.Ctx) error {
	var prevData map[string]interface{}

	err := json.Unmarshal(c.Request().Body(), &prevData)

	if err != nil {
		panic("invalid data")
	}
	fmt.Println(prevData)

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	prevData["User"] = claims

	js, _ := json.Marshal(prevData)
	c.Request().SetBody([]byte(js))
	return c.Next()
}

func restricted(c *fiber.Ctx) error {
	fmt.Println(string(c.Request().Body()))
	var data interface{}
	err := json.Unmarshal(c.Request().Body(), &data)
	if err != nil {
		panic("Invalid json")
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// js, _ := json.Marshal(claims)
	return c.JSON(claims)
}
