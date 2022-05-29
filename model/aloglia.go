package model

import (
	"fmt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"os"
)

const (
	ALGOLIA_APP_ID     = "4AXYXC3HDX"
	ALGOLIA_INDEX_NAME = "kefu"
)

var algolia *search.Index

type QA struct {
	ID       string `json:"objectID"`
	Question string `json:"q"`
	Answer   string `json:"a"`
}

func init() {
	appID := ALGOLIA_APP_ID
	ALGOLIA_API_KEY := os.Getenv("ALGOLIA_API_KEY")
	apiKey := ALGOLIA_API_KEY
	indexName := ALGOLIA_INDEX_NAME

	algolia = search.NewClient(appID, apiKey).InitIndex(indexName)
}

func searchAlgolia(question string) *[]QA {
	res, err := algolia.Search(question)
	if err != nil {
		fmt.Println(err)
	}

	var contacts []QA

	err = res.UnmarshalHits(&contacts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("search results: ", contacts)
	return &contacts
}
