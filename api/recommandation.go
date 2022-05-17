package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

func Recommandation(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json; application/json")
	if err := r.ParseForm(); err != nil {
		err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
		log.Println(err)
		http.Error(w, err_msg, http.StatusInternalServerError)
		fmt.Fprintf(w, "")
		return
	}

	res := domain.NewDummyResults(10)
	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "")
	}
	_, err = io.WriteString(w, string(b))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Fprintf(w, "")
	}

}
