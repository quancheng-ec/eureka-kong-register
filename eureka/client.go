package eureka

import (
	"log"
	"regexp"
	"sync"

	"github.com/hudl/fargo"
)

type Client struct {
	fargo.EurekaConnection
}

type Config struct {
	ServiceUrls         []string
	PollIntervalSeconds int
}

func NewClient(c Config) Client {
	var eurekaConf struct {
		InTheCloud            bool
		ConnectTimeoutSeconds int
		UseDNSForServiceUrls  bool
		DNSDiscoveryZone      string
		ServerDNSName         string
		ServiceUrls           []string
		ServerPort            int
		ServerURLBase         string
		PollIntervalSeconds   int
		EnableDelta           bool
		PreferSameZone        bool
		RegisterWithEureka    bool
		Retries               int
	}

	eurekaConf.ServiceUrls = c.ServiceUrls
	eurekaConf.PollIntervalSeconds = c.PollIntervalSeconds
	eurekaConf.Retries = 3
	eurekaConf.ServerURLBase = "eureka/v2"

	eurekaClient := fargo.NewConnFromConfig(fargo.Config{
		Eureka: eurekaConf,
	})

	return Client{
		eurekaClient,
	}
}

func (client *Client) StartEurekaPolling(onPoll func(app *fargo.Application), filter string) {

	registeredApps, _ := client.GetApps()

	if len(filter) == 0 {
		filter = ".*"
	}

	filterReg := regexp.MustCompile(filter)

	var waitGroup sync.WaitGroup

	for _, app := range registeredApps {

		waitGroup.Add(1)

		if filterReg.MatchString(app.Name) {
			c := make(chan struct{})
			go func(app *fargo.Application) {
				for update := range client.ScheduleAppUpdates(app.Name, true, c) {
					if update.Err != nil {
						log.Printf("Most recent request for application %q failed: %v\n", app.Name, update.Err)
						waitGroup.Done()
						continue
					}
					onPoll(update.App)

				}
			}(app)
		}
	}

	waitGroup.Wait()
}
