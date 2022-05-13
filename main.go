package main

// #cgo pkg-config: python3
// #include <Python.h>
import "C"

import (
	"fmt"
	"net/http"
	"os"

	"github.com/erudit-recommandation/search-engine-webapp/server"
)

func main() {

	fmt.Printf("\nServer started at: http://localhost%v", server.GetPort())
	srv := &http.Server{}
	close := make(chan os.Signal, 1)
	server.BuildServer(srv, close)

}
