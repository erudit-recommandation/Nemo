package middleware

import (
	"errors"
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

var PERSONA_IMAGE_DIRECTORY_TEMPLATE = "./static/images/persona/%v"
var PERSONA_IMAGE_TEMPLATE = "./static/images/persona/%v/%v.svg"

type cache map[uint32]cacheElement[any]

func (c *cache) ClearExpired() {
	keys := reflect.ValueOf(*c).MapKeys()
	for _, k := range keys {
		if (*c)[uint32(k.Uint())].IsExpired() {
			log.Printf("cache %v removed", k.Uint())
			path := fmt.Sprintf(PERSONA_IMAGE_DIRECTORY_TEMPLATE, uint32(k.Uint()))
			os.RemoveAll(path)
			delete(*c, uint32(k.Uint()))
		}
	}
}

type cacheElement[T any] struct {
	CreatedDate time.Time
	Query       string
	Elements    []T
	PageSize    uint
}

func newCacheElement(query string, hasedQuery uint32, elements []any, pageSize uint) cacheElement[any] {
	createPersonaDirectory(hasedQuery)
	return cacheElement[any]{
		Query:       query,
		CreatedDate: time.Now(),
		Elements:    elements,
		PageSize:    pageSize,
	}
}
func (c cacheElement[T]) IsExpired() bool {
	return -1*time.Until(c.CreatedDate) >= CACHE_DURATION
}

func (c cacheElement[T]) NumberOfPage() int {
	n := int(math.Ceil(float64(len(c.Elements))/float64(c.PageSize))) - 1
	if n <= 0 {
		n = -1
	}
	return n
}

func (c cacheElement[T]) GetPage(page uint) ([]T, error) {
	if c.NumberOfPage() == -1 {
		return c.Elements[c.PageSize*page:], nil
	}

	if int(page) > c.NumberOfPage() {
		return nil, fmt.Errorf("la page n'existe pas")
	}

	if int(page) == c.NumberOfPage() {
		return c.Elements[c.PageSize*page:], nil
	}

	return c.Elements[c.PageSize*page : c.PageSize*(page+1)], nil
}

func createPersonaDirectory(hasedQuery uint32) error {
	path := fmt.Sprintf(PERSONA_IMAGE_DIRECTORY_TEMPLATE, hasedQuery)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func createPersonaSVG(articles []domain.Article, hasedQuery uint32) error {
	if len(articles) == 0 {
		return nil
	}
	for _, a := range articles {
		if a.PersonaSvg != "" {
			path := fmt.Sprintf(PERSONA_IMAGE_TEMPLATE, hasedQuery, a.ID)
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				f, err := os.Create(path)
				if err != nil {
					return err
				}
				b := []byte(a.PersonaSvg)
				if _, err := f.Write(b); err != nil {
					return err
				}
				f.Close()
			}
		}

	}

	return nil
}

func GetArticleFromCache(hasedQuery uint32, limit uint, currentCache *cache) ([]domain.Article, error) {
	resp := make([]domain.Article, 0, limit)
	cacheValue := (*currentCache)[hasedQuery]

	articles, err := cacheValue.GetPage(0)

	if err != nil {
		return nil, err
	}

	for _, a := range articles {
		resp = append(resp, a.(domain.Article))
	}

	return resp, nil
}

func hash(s string, corpus string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s + "_" + corpus))
	return h.Sum32()
}
