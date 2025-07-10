package main

import (
	"context"

	"github.com/pkossyfas/go-server-bootstrap/config"
	db "github.com/pkossyfas/go-server-bootstrap/dbconnector"
	"github.com/pkossyfas/go-server-bootstrap/logger"
	"github.com/pkossyfas/go-server-bootstrap/server"
)

func main() {
	// oad app configuration
	config.LoadAppConfig()

	// Init db connection.
	err := db.InitDBConn(
		context.Background(),
		config.GetConfig.DBHost,
		config.GetConfig.DBPort,
		config.GetConfig.DBUser,
		config.GetConfig.DBPassword,
		config.GetConfig.DBName,
	)
	if err != nil {
		// In case creating the db connection fails then the app should exit,
		// but for this sample app we just log it as an error in order to
		// demonstrate the functioanlity of the ready endpoint.
		logger.Error(err, "error creating the db connection")
	}

	// Load endpoint and start the server.
	server.LoadEndpoints()
	server.StartServer()
}
