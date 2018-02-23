package main

import (
	"fmt"

	"github.com/hudl/fargo"
)

func main() {

	eurekaClient := fargo.NewConn("http://eureka.dev.quancheng-ec.com/eureka")

	registeredApps, err := eurekaClient.GetApps()

	if err != nil {
		fmt.Println("fetch error:", err)
	}

	for _, app := range registeredApps {
		fmt.Println("name",app.Name)

	}

}
