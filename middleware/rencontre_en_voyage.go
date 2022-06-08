package middleware

import (
	"net/http"
	"os/exec"
)

var gemsim_service = exec.Command("python3", "./erudit_notebooks/gemsim_service.py")

func RencontreEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		next(w, req)
	}
}

func SetGemsimServicePath(path string) {
	gemsim_service = exec.Command("python3", path)
}
