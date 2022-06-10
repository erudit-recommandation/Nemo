package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

var gemsim_service *exec.Cmd = exec.Command("python3", "./scripts/gemsim_service.py")

//var gemsim_service *exec.Cmd = exec.Command("echo", "-n", `{"Name": 11.1, "Age": 32.2}`)
var gemsim_service_listener io.ReadCloser = nil
var gemsim_service_writer io.WriteCloser = nil
var gemsim_is_on bool = false

func RencontreEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	if !gemsim_is_on {
		log.Println("\n\nstarting gemsim service")
		start_gemsim_service()
	}
	return func(w http.ResponseWriter, req *http.Request) {
		query := "test"
		io.WriteString(gemsim_service_writer, query)
		log.Printf("query: %v", query)
		var res map[string]float64

		if err := json.NewDecoder(gemsim_service_listener).Decode(&res); err != nil {
			log.Fatal(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}
		log.Println(res)
		next(w, req)
	}
}

func SetGemsimServicePath(path string) {
	gemsim_service = exec.Command("python3", path)
	start_gemsim_service()
}

func start_gemsim_service() {
	stdout, err := gemsim_service.StdoutPipe()
	gemsim_service_listener = stdout
	if err != nil {
		error_msg := fmt.Sprintf("was not able to start the gemsim service with error \n%v", err)
		log.Println(error_msg)
		panic(error_msg)
	}

	gemsim_service_writer, err = gemsim_service.StdinPipe()
	if err != nil {
		error_msg := fmt.Sprintf("was not able to start the gemsim service with error \n%v", err)
		log.Println(error_msg)
		panic(error_msg)
	}

	go func() {
		err = gemsim_service.Run()
		log.Println("gemsim service started")
		if err != nil {
			error_msg := fmt.Sprintf("was not able to start the gemsim service with error \n%v", err)
			log.Println(error_msg)
			panic(error_msg)
		}
		log.Println("gemsim service exited")
	}()
	gemsim_is_on = true

}
