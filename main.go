package main

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/op/go-logging"
	"github.com/zephyrpersonal/eureka-kong-register/app"
)

func main() {
	logging.SetLevel(logging.ERROR, "fargo")

	config := app.Config{}

	err := env.Parse(&config)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	fmt.Printf("%+v\n", config)

	c := app.NewApp(config)

	c.Start()

}
