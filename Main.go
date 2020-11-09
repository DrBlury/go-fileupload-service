package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
)

func main() {
	configFiber := new(fiber.Config)
	configFiber.BodyLimit = 700 * 1024 * 1024

	app := fiber.New(*configFiber)

	// Load defaults.yml
	var defaults conf
	defaults.getDefaults()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Banana 👋!")
	})

	app.Post("/upload/generic", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["documents"]
			for _, file := range files {
				fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
				c.SaveFile(file, fmt.Sprintf("./tempFiles/%s", file.Filename))
			}
		}
		return c.SendString("uploaded")
	})

	app.Post("/upload/picture", func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			files := form.File["documents"]
			for _, file := range files {
				datatype := file.Header["Content-Type"][0]
				substrings := strings.Split(datatype, "/")

				if substrings[0] != "image" {
					return c.SendString("Please upload an image")
				}

				c.SaveFile(file, fmt.Sprintf("./pictures/%s", file.Filename))
			}
		}
		return c.SendString("uploaded")
	})

	app.Static("/", "./public")
	app.Static("/picture", "./pictures")
	app.Static("/files", "./tempFiles")

	// Set default port if not set
	port, isPresent := os.LookupEnv("fiber_port")
	if isPresent == false {
		port = defaults.Port
	}

	app.Listen(":" + port)
}
