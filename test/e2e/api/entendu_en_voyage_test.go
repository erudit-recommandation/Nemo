package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/api"
	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func TestMain(m *testing.M) {
	env := config.EnvVariable{
		ArangoPort:              "http://localhost:8529",
		ArangoPassword:          "rootpassword",
		ArangoUsername:          "root",
		ArangoDatabase:          "erudit",
		ArangoArticleCollection: "articles",
	}

	config.SetConfig(&env)

	os.Exit(m.Run())
}

func TestAPIEntenduEnVoyage(t *testing.T) {

	data := url.Values{}
	data.Set("text", "la crise des logements")

	req, err := http.NewRequest(http.MethodPost, "/api/entendu_en_voyage", strings.NewReader(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.EntenduEnVoyage(api.EntenduEnVoyage))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v with the error %v",
			status, http.StatusOK, rr.Body.String())
	}
	var articles []domain.Article

	err = json.Unmarshal(rr.Body.Bytes(), &articles)
	if err != nil {
		t.Fatal(err)
	}

	if len(articles) == 0 {
		t.Error("their should be some recommandations")
	}
}

func TestAPIEntenduEnVoyageNoText(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "/api/entendu_en_voyage", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.EntenduEnVoyage(api.EntenduEnVoyage))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v, shoud return an error because their is no text sended",
			status, http.StatusInternalServerError)
	}

}
