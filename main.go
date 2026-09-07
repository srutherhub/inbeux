package main

import (
	"context"
	msg "inbeux/internal/message"
	"inbeux/internal/platform"
	"inbeux/internal/platform/db"
	"inbeux/internal/user"
	"log"
	"os"
	"os/signal"
	"syscall"

	c "github.com/srutherhub/web-app/controller"
	s "github.com/srutherhub/web-app/server"
)

func main() {
	serverConfig := s.InitServerCfg("8080")
	server := s.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	dbConn, close, err := db.InitializeDb(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.RunMigrations()

	defer close()

	baseController := baseController()
	server.RegisterController(*baseController)

	userRepository := user.NewRepository(dbConn)
	userService := user.NewService(userRepository)

	messageRepository := msg.NewRepository(dbConn)
	messageService := msg.NewService(messageRepository)
	messageController := messageController(messageService, userService)
	server.RegisterController(*messageController)

	go messageService.StartPendingMessagesBatch(ctx)

	go func() {
		server.Start(serverConfig)
	}()

	<-ctx.Done()
}

func messageController(es *msg.MessageService, us msg.UserProvider) *c.Controller {

	messageController := c.New()
	messageController.SetBase("/email")
	messageController.RegisterRoute(c.Route{Method: "POST", Path: "/receiver", Handler: msg.EmailReceiverHandler(es, us)})

	return messageController
}

func baseController() *c.Controller {
	baseController := c.New()
	baseController.SetBase("")
	baseController.RegisterRoute(c.Route{Method: "GET", Path: "/healthz", Handler: platform.HealthzHandler()})

	return baseController
}
