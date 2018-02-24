package main

import (
	"eureka-kong-register/app"

	"github.com/op/go-logging"
)

func main() {
	logging.SetLevel(logging.ERROR, "fargo")

	c := app.NewApp(app.Config{
		KongHost:     "http://qccost-gateway-admin.dev.quancheng-ec.com",
		PollInterval: 10,
	})

	c.Start()

}
