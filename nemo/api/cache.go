package api

import (
	"encoding/json"
	"log"
	"net/http"

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

func DeleteCache(w http.ResponseWriter, r *http.Request) {
	middleware.CACHE_ENTENDU_EN_VOYAGE.Clear()
	middleware.CACHE_RENCONTRE_EN_VOYAGE.Clear()
	middleware.CACHE_ACCOSTE_EN_VOYAGE.Clear()
}
