package main

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"
	"urbangrid.com/constants"
	"urbangrid.com/database"
	"urbangrid.com/handlers"
	"urbangrid.com/queues"
	"urbangrid.com/workers"
)

func main() {
	defer database.DisconnectMongo()
	defer database.DisconnectRedis()
	defer queues.DisconnectQueue()

	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		Format:     "[${ip}]:${port} ${time} ${status} - ${method} ${latency} ${path}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
		TimeZone:   "Asia/Kolkata",
	}))

	accidents := app.Group("/accidents")
	{
		accidents.Get("/", handlers.GetAccidents)
		accidents.Post("/", handlers.UploadAccident)

		accidents.Static("/", "./uploads")
	}

	signales := app.Group("/signales")
	{
		signales.Get("/", handlers.GetSignalUpdate)
		signales.Post("/", handlers.SetSignal)
	}

	ws := app.Group("/ws")
	{
		ws.Get("/signal", websocket.New(handlers.SignalSocket, websocket.Config{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		}))

		ws.Get("/accident", websocket.New(handlers.AccidentSocket, websocket.Config{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
		}))
	}

	go func() {
		if err := app.Listen(":8726"); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Initialize and start the Asynq server in a separate goroutine
	go func() {
		mux := asynq.NewServeMux()
		mux.HandleFunc(constants.ACCIDENT_SMS, workers.SendSMS)
		mux.HandleFunc(constants.SIGNAL_CHANGE, workers.SignalChangeNotification)

		if err := queues.Server.Run(mux); err != nil {
			log.Fatalf("Failed to run Asynq server: %v", err)
		}
	}()

	// Initialize and start the cron job
	go func() {
		c := cron.New()

		// Schedule the update function to run every minute
		c.AddFunc("@every 2s", workers.TrafficController)

		c.Start()

		select {}
	}()

	select {}
}
