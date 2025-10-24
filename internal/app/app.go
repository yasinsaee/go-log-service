package app

import "github.com/yasinsaee/go-log-service/internal/app/config"

func StartApp() {
	config.LoadEnv()

	InitElastic()

	StartGRPCServer()

}
