package main

import (
	"ticket_booking_app/handlers"
	"ticket_booking_app/repository"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Ticket Booking App",
		ServerHeader: "Fiber",
	})

	// repository
	eventRepository := repository.NewEventRepository(nil)

	// routing
	server := app.Group("/api")

	// handlers
	handlers.NewEventHandler(server.Group("/event"),eventRepository)

	app.Listen(":3000")
}