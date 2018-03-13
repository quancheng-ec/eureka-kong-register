package app

import (
	"log"
	"time"

	"github.com/zephyrpersonal/eureka-kong-register/eureka"
	"github.com/zephyrpersonal/eureka-kong-register/kong"
)

type Config struct {
	KongHost     string   `env:"KONG_HOST" envDefault:"http://qccost-gateway-admin.dev.quancheng-ec.com"`
	EurekaUrls   []string `env:"EUREKA_URLS" envSeparator:"|" envDefault:"http://eureka.dev.quancheng-ec.com/eureka"`
	PollInterval int      `env:"POLL_INTERVAL" envDefault:"10"`
	Filter       string   `env:"FILTER" envDefault:".*"`
}

type App struct {
	config       Config
	KongClient   kong.Client
	EurekaClient eureka.Client
}

func NewApp(c Config) (app App) {
	return App{
		config: c,
		KongClient: kong.NewClient(kong.Config{
			Host: c.KongHost,
		}),
		EurekaClient: eureka.NewClient(eureka.Config{
			ServiceUrls:         c.EurekaUrls,
			PollIntervalSeconds: c.PollInterval,
		}),
	}
}

func (a *App) Start() {
	for {
		log.Println("poll start")
		a.EurekaClient.StartEurekaPollingFetch(a.KongClient.RegisterUpstream, a.config.Filter)
		time.Sleep(time.Second * time.Duration(a.config.PollInterval))
	}
}
