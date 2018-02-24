package app

import (
	"eureka-kong-register/eureka"
	"eureka-kong-register/kong"

	"github.com/hudl/fargo"
)

type Config struct {
	KongHost     string
	PollInterval int
}

type App struct {
	config       Config
	KongClient   kong.Client
	EurekaClient fargo.EurekaConnection
}

func NewApp(c Config) (app App) {
	return App{
		config: c,
	}
}

func (a *App) init() {
	a.KongClient = kong.NewClient(kong.Config{
		Host: a.config.KongHost,
	})
}

func (a *App) Start() {
	a.init()
	eureka.StartEurekaPolling("eureka-config.gcfg", a.config.PollInterval, a.KongClient.RegisterUpstream, "NODE")
}
