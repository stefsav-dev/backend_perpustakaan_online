package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"backend_perpustakaan_online/config"
	"backend_perpustakaan_online/database"
	"backend_perpustakaan_online/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	config.ConnectDB()

	database.Migrate()

	if os.Getenv("SEED_DATA") == "true" {
		database.Seeder()
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Library API",
		ErrorHandler: errorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	bookHandler := handlers.NewBookHandler()

	api := app.Group("/api")

	books := api.Group("/books")
	books.Get("/", bookHandler.GetAllBooks)
	books.Get("/:id", bookHandler.GetBookByID)
	books.Post("/", bookHandler.CreateBook)
	books.Put("/:id", bookHandler.UpdateBook)
	books.Patch("/:id", bookHandler.UpdateBook)
	books.Delete("/:id", bookHandler.DeleteBook)
	books.Patch("/:id/status", bookHandler.UpdateBookStatus)

	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := config.DB.DB()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection error",
			})
		}

		if err := sqlDB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Database ping failed",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "ok",
			"message":  "Library API is running",
			"database": "connected",
		})
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Endpoint not found",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func errorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"success": false,
		"error":   "Internal server error",
	})
}
