package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/barkhayot/auth-fiber/models"
	"github.com/barkhayot/auth-fiber/database"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"
)

const SecretKey = "secret"

// Register Function 
func Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parsing the data that has been entered by User
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Bcrypting the password
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name  : data["name"],
		Email : data["email"],
		Password : password, 
	}

	database.DB.Create(&user)

	return c.JSON(user)
}

// Login Function 
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parsing the data that has been entered by User
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}
	
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "couldn't login",
		})
	}

	cookie := fiber.Cookie{
		Name:		"jwt",
		Value :		token,
		Expires: 	time.Now().Add(time.Hour * 24),
		HTTPOnly: 	true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message" : "Success",
	})
}

// User  Function
func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

// Logout Function
func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
