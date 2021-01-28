package algolia

import (
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"log"
)

type Algolia struct {
	client *search.Client
	index  *search.Index
}

func New(appID, apiKey string) *Algolia {
	client := search.NewClient(appID, apiKey)
	return &Algolia{
		client: client,
	}
}

func (a *Algolia) InitializeIndex(name string) {
	log.Println("Index Name: ", name)
	a.index = a.client.InitIndex(name)

}

func (a *Algolia) CreateMany(v interface{}) error {
	_, err := a.index.SaveObjects(v)
	if err != nil {
		return err
	}
	return nil
}

func (a *Algolia) UpdateMany(v interface{}) error {
	_, err := a.index.SaveObjects(v)
	return err
}

func (a *Algolia) DeleteMany(oid []string) error {
	_, err := a.index.DeleteObjects(oid)
	return err
}
