package routes

import (
	"strings"

	"github.com/acatalepsy17/pigeon/models"
	"github.com/acatalepsy17/pigeon/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUser(token string, db *gorm.DB) (*models.User, *string) {
	if !strings.HasPrefix(token, "Bearer ") {
		err := "Bearer token is not provided!"
		return nil, &err
	}
	user, err := DecodeAccessToken(token[7:], db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (ep Endpoint) AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB

	if len(token) < 1 {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_UNAUTHORIZED_USER, "Unauthorized User!"))
	}
	user, err := GetUser(token, db)
	if err != nil {
		return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
	}
	c.Locals("user", user)
	return c.Next()
}

func (ep Endpoint) GuestMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := ep.DB
	var user *models.User
	if len(token) > 0 {
		userObj, err := GetUser(token, db)
		if err != nil {
			return c.Status(401).JSON(utils.RequestErr(utils.ERR_INVALID_TOKEN, *err))
		}
		user = userObj
	}
	c.Locals("user", user)
	return c.Next()
}
