package scrap

import (
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// providers
const (
	AmazonPrime = "AmazonPrime"
)

var validators = map[string]*regexp.Regexp{
	AmazonPrime: regexp.MustCompile("[A-Z0-9]{8,14}"),
}

// IsValidID validates an id for a provider
// TODO: returning true if the provider is not found is fine for the scope of the project but
// but otherwise should be considered more carefully
func IsValidID(provider, id string) bool {
	r, found := validators[provider]
	if !found {
		log.Debugln("validator regexp for provider %v not found", provider)
		return true
	}
	return r.MatchString(id)
}

// the queries map contains mapping: provider -> query
// this should be probably written on a fql file and bundled
// or distributed separately
var queries = map[string]string{
	AmazonPrime: `
LET doc = DOCUMENT("https://www.amazon.%s/gp/product/%s")
LET title = ELEMENT(doc, '[data-automation-id="title"]')
LET year = ELEMENT(doc, '[data-automation-id="release-year-badge"]')
// image
// LET packshot = ELEMENT(doc, '.dv-fallback-packshot-image') 
// LET img = ELEMENT(packshot, 'img')

// actors
LET meta = ELEMENT(doc, '[data-automation-id="meta-info"]')
LET actors = (
	FOR a IN ELEMENTS(meta, 'a')
		// this is just ridiculous
        LET urlc = SPLIT(a.attributes.href, "=")
    	FILTER urlc[1] == "atv_dp_pd_star?phrase"
    	RETURN TRIM(a.innerHTML)
)

// suggested
LET recommendation = ELEMENT(doc, '.dv-hover-packshot-container')
LET similar = (
	FOR a IN ELEMENTS(recommendation, 'a')
		// this one also should be extracted with REGEXP (but currently does not work properly)
		RETURN SPLIT(a.attributes.href, "/")[4]
)

RETURN {
  title: TRIM(title.innerText),
  release_year: TRIM(year.innerText),
  actors: actors,
  similar_ids: similar,
  //poster: img.attributes.src,
} `}

// this should be expanded to accept other parmaterers like language
func buildQuery(provider, id, locale string) (q string, err error) {
	q, found := queries[provider]
	if !found {
		err = fmt.Errorf("provider not supported: %s", provider)
		return
	}
	if !IsValidID(provider, id) {
		err = fmt.Errorf("invalid id %v for provider %v", id, provider)
		return
	}
	q = fmt.Sprintf(q, locale, id)
	return
}
