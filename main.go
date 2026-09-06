package main

import (
	"context"
	msg "inbeux/internal/message"
	"inbeux/internal/platform/db"
	"inbeux/internal/user"
	"log"

	c "github.com/srutherhub/web-app/controller"
	s "github.com/srutherhub/web-app/server"
)

func main() {
	serverConfig := s.InitServerCfg("8080")
	server := s.New()

	ctx := context.Background()
	dbConn, close, err := db.InitializeDb(ctx)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.RunMigrations()

	defer close()

	userRepository := user.NewRepository(dbConn)
	userService := user.NewService(userRepository)

	messageRepository := msg.NewRepository(dbConn)
	messageService := msg.NewService(messageRepository)
	messageController := messageController(messageService, userService)
	server.RegisterController(*messageController)

	server.Start(serverConfig)
}

func messageController(es *msg.MessageService, us *user.UserService) *c.Controller {

	messageController := c.New()
	messageController.SetBase("/email")
	messageController.RegisterRoute(c.Route{Method: "POST", Path: "/receiver", Handler: msg.EmailReceiverHandler(es, us)})

	return messageController
}
