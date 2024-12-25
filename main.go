package main

import (
	"log"
	"runtime"

	"github.com/acatalepsy17/yappy/config"
	"github.com/acatalepsy17/yappy/database"
	"github.com/acatalepsy17/yappy/routes"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.GetConfig()
	db := database.ConnectDb(cfg)
	sqlDb, _ := db.DB()
	app := fiber.New(fiber.Config{
		Concurrency: runtime.NumCPU(),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	swaggerCfg := swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "/",
		Title:    "Yappy API Specification",
		CacheAge: 1,
	}

	app.Use(swagger.New(swaggerCfg))
	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	routes.SetupRoutes(app, db)
	defer sqlDb.Close()
	log.Fatal(app.Listen("127.0.0.1:8000"))
}