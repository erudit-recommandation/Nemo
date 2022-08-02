package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/server"
)

func main() {

	devModeDocker := flag.Bool("dd", false, "developpement model")
	devMode := flag.Bool("d", false, "developpement model")
	productionMode := flag.Bool("p", false, "developpement model")
	flag.Parse()

	if *devModeDocker {
		config.CONFIG_MODE = config.DEV_DOCKER
		fmt.Println("developpement mode with docker")
	} else if *devMode {
		config.CONFIG_MODE = config.DEV
		fmt.Println("developpement mode")
	} else if *productionMode {
		config.CONFIG_MODE = config.PRODUCTION
		fmt.Println("production mode")
	} else {
		fmt.Println("developpement mode")
		config.CONFIG_MODE = config.DEV
	}

	fmt.Printf("\nServer started at: http://localhost%v", server.GetPort())
	srv := &http.Server{}
	close := make(chan os.Signal, 1)
	server.BuildServer(srv, close)

}
