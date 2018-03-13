package eureka

import (
	"regexp"

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

func (client *Client) StartEurekaPollingFetch(onPoll func(app *fargo.Application), filter string) {
	registeredApps, _ := client.GetApps()

	if len(filter) == 0 {
		filter = ".*"
	}

	filterReg := regexp.MustCompile(filter)

	ch := make(chan string)

	for _, app := range registeredApps {
		if filterReg.MatchString(app.Name) {
			go func(app *fargo.Application) {
				onPoll(app)
				ch <- app.Name
			}(app)

			<-ch
		}
	}

	close(ch)
}
