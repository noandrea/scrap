package scrap

import (
	"context"
	"encoding/json"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	log "github.com/sirupsen/logrus"
)

// Movie is a movie to scrape
type Movie struct {
	Title       string   `json:"title"`
	ReleaseYear string   `json:"release_year"`
	Actors      []string `json:"actors"`
	Poster      string   `json:"poster"`
	SimilarIds  []string `json:"similar_ids"`
}

var (
	ferret *compiler.Compiler
	ctx    context.Context
)

// Configure scrap
func Configure(chromeAddress string) {
	ferret = compiler.New()
	// create a root context
	ctx = context.Background()
	// chrome headless is required to execute js
	ctx = drivers.WithContext(ctx, cdp.NewDriver(
		cdp.WithAddress(chromeAddress)),
		drivers.AsDefault(),
	)
}

// Run execute the data extraction
func Run(provider, id, region string, m interface{}) (err error) {
	log.Debugf("extraction provider:%s id:%s region:%s", provider, id, region)
	// retrieve the query for the provider
	query, err := buildQuery(provider, id, region)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("query is %s", query)
	// compile the query
	program, err := ferret.Compile(query)
	if err != nil {
		log.Errorf("error compiling query: %v", err)
		return
	}
	// execute the query
	out, err := program.Run(ctx)
	if err != nil {
		log.Errorf("error executing query: %v", err)
		return
	}
	// unmarshal result
	err = json.Unmarshal(out, &m)
	if err != nil {
		log.Errorf("error unmarshal result: %v", err)
		return
	}
	return
}
