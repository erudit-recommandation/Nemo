package middleware

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

type cache map[uint32]cacheElement

func (c *cache) ClearExpired() {
	keys := reflect.ValueOf(*c).MapKeys()
	for _, k := range keys {
		if (*c)[uint32(k.Uint())].IsExpired() {
			log.Printf("cache %v removed", k.Uint())
			path := fmt.Sprintf("./static/images/persona/%v", uint32(k.Uint()))
			os.RemoveAll(path)
			delete(*c, uint32(k.Uint()))
		} else {
			log.Println(time.Until((*c)[uint32(k.Uint())].CreatedDate))
		}
	}
}

type cacheElement struct {
	CreatedDate time.Time
	Query       string
	Elements    []infrastructure.ArticlesID
}

func newCacheElement(query string, hasedQuery uint32, elements []infrastructure.ArticlesID) cacheElement {
	createPersonaDirectory(hasedQuery)
	return cacheElement{
		Query:       query,
		CreatedDate: time.Now(),
		Elements:    elements,
	}
}
func (c cacheElement) IsExpired() bool {
	return -1*time.Until(c.CreatedDate) >= CACHE_DURATION
}

func (c cacheElement) NumberOfPage() uint {
	return uint(math.Ceil(float64(len(c.Elements)) / float64(MAX_PAGE_ENTENDU_EN_VOYAGE)))
}

func (c cacheElement) GetPage(page uint) ([]infrastructure.ArticlesID, error) {
	if page > c.NumberOfPage() {
		return nil, fmt.Errorf("la page n'existe pas")
	}

	if page == c.NumberOfPage() {
		return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page:], nil
	}

	return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page : MAX_PAGE_ENTENDU_EN_VOYAGE*(page+1)], nil
}

func createPersonaDirectory(hasedQuery uint32) error {
	path := fmt.Sprintf("./static/images/persona/%v", hasedQuery)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func createPersonaSVG(article []domain.Article) {

}
