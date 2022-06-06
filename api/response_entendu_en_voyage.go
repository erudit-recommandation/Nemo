package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func EntenduEnVoyage(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json; application/json")

	_, err := io.Copy(w, r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "")
	}

}
