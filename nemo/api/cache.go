package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func GetCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	caches := map[string]middleware.Cache{
		"entendu_en_voyage":   middleware.CACHE_ENTENDU_EN_VOYAGE,
		"rencontre_en_voyage": middleware.CACHE_RENCONTRE_EN_VOYAGE,
		"accoster_en_voyage":  middleware.CACHE_ACCOSTE_EN_VOYAGE,
	}
	jsonResp, err := json.Marshal(caches)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)

}

var last_time_deleted = time.Now()

var MIN_DURATION_DELETION = 5 * time.Minute

func DeleteCache(w http.ResponseWriter, r *http.Request) {
	if -1*time.Until(last_time_deleted) >= MIN_DURATION_DELETION {
		middleware.CACHE_ENTENDU_EN_VOYAGE.Clear()
		middleware.CACHE_RENCONTRE_EN_VOYAGE.Clear()
		middleware.CACHE_ACCOSTE_EN_VOYAGE.Clear()
		last_time_deleted = time.Now()

	} else {
		log.Println(1 * time.Until(last_time_deleted))
		http.Error(w, fmt.Sprintf("Le cache peut être supprimé dans %v", MIN_DURATION_DELETION+time.Until(last_time_deleted)), http.StatusForbidden)
	}

}
