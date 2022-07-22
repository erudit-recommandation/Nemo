package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func JSONResult(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json; application/json")

	var resp middleware.ResultResponse

	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	articles := resp.Data

	b, err := json.Marshal(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.WriteString(w, string(b))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "")
	}

}
