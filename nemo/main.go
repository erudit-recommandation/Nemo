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

	devMode := flag.Bool("d", false, "developpement model")
	flag.Parse()
	fmt.Println(*devMode)

	config.DEV_MODE = *devMode
	if config.DEV_MODE {
		fmt.Println("developpement mode")
	} else {
		fmt.Println("production mode")
	}

	fmt.Printf("\nServer started at: http://localhost%v", server.GetPort())
	srv := &http.Server{}
	close := make(chan os.Signal, 1)
	server.BuildServer(srv, close)

}
