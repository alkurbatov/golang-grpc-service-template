// Package main is the main entrypoint of the service.
package main

import (
	"os"

	"github.com/alkurbatov/golang-grpc-service-template/internal/app"
	"github.com/alkurbatov/golang-grpc-service-template/internal/config"
)

func main() {
	cfg := config.New()

	if err := app.Run(cfg); err != nil {
		os.Exit(1)
	}
}
