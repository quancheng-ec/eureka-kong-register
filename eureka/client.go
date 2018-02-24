package eureka

import (
	"log"
	"regexp"

	"github.com/hudl/fargo"
)

func StartEurekaPolling(location string, pollInterval int, onPoll func(app *fargo.Application), filter string) {
	client, _ := fargo.NewConnFromConfigFile(location)

	registeredApps, _ := client.GetApps()

	if len(filter) == 0 {
		filter = ".*"
	}

	filterReg := regexp.MustCompile(filter)

	for _, app := range registeredApps {

		if filterReg.MatchString(app.Name) {
			c := make(chan struct{})
			for update := range client.ScheduleAppUpdates(app.Name, true, c) {
				if update.Err != nil {
					log.Printf("Most recent request for application %q failed: %v\n", app.Name, update.Err)
					continue
				}
				//onPoll(update.App)
				log.Printf("Application %q has %d instances.\n", app.Name, len(update.App.Instances))
			}

			log.Println(1)

			close(c)
		}
	}

}
