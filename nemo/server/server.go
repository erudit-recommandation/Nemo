package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/api"
	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/middleware"
	"github.com/erudit-recommandation/search-engine-webapp/route"
	"github.com/gorilla/mux"
)

func setRoute(r *mux.Router) {
	r.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		fmt.Fprintf(w, "<h1>Boom Boom</h1>")
	}).Methods("GET")

	r.HandleFunc("/", route.Homepage).Methods("GET")
	r.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		middleware.Error(w, r, 200, "Aucune erreur, voici la page test")
	}).Methods("GET")

	r.HandleFunc(route.ENTENDU_EN_VOYAGE, middleware.EntenduEnVoyage(route.Result)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("%v/{hashedQuery}", route.ENTENDU_EN_VOYAGE), middleware.EntenduEnVoyageCached(route.Result)).Methods("GET")

	r.HandleFunc(route.RENCONTRE_EN_VOYAGE, middleware.RencontreEnVoyage(route.Result)).Methods("POST")

	r.HandleFunc(route.ACCOSTE_EN_VOYAGE, middleware.AccosteEnVoyage(route.Result)).Methods("GET")

	r.HandleFunc(route.CLEF_CANONIQUE, func(w http.ResponseWriter, r *http.Request) {
		middleware.Error(w, r, 501, "Service à venir!")
	}).Methods("POST")

	r.HandleFunc(route.VAISEAU, route.Vaiseau).Methods("GET")

	r.HandleFunc(route.REMERCIEMENTS, route.Remerciements).Methods("GET")

	r.HandleFunc("/cache", api.GetCache).Methods("GET")
	r.HandleFunc("/cache", api.DeleteCache).Methods("DELETE")

	r.NotFoundHandler = r.NewRoute().HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.Error(w, r, http.StatusNotFound, "page non trouvée")
	}).GetHandler()

}

func GetPort() string {

	return ":" + config.GetConfig().Port
}

func BuildServer(srv *http.Server, close chan os.Signal) {

	r := initializeServer()
	setRoute(r)

	srv.Addr = GetPort()
	srv.Handler = r

	signal.Notify(close, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	<-close
	closeServer(srv)
}

// initialise the router
func initializeServer() *mux.Router {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))
	static := r.PathPrefix("/static/")
	fs := http.FileServer(http.Dir("static/"))

	static.Handler(http.StripPrefix("/static/", fs))
	return r
}

// send a close server signal
func SendCloseSignal(close chan os.Signal) {
	close <- os.Interrupt
}

func closeServer(srv *http.Server) {
	log.Println("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
