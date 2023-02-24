package main

import (
	"log"

	"github.com/RipperAcskt/innotaxi/internal/app"
)

// @title InnoTaxi API
// @version 1.0
// @description API for order taxi
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  ripper@gmail.com

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
