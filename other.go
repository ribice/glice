package main

// Resolve indirect repos as described here:
// https://golang.org/cmd/go/#hdr-Remote_import_paths

import (
	"fmt"
	"net/http"
	"strings"

	m "github.com/keighl/metabolize"
	"github.com/ribice/glice/api"
)

// MetaData contains tags for fetching meta tags
type MetaData struct {
	Import string `meta:"go-import"`
	Source string `meta:"go-source"`
}

var cache = map[string]*api.Repository{}

func getOtherRepo(ptrDep *string, verbose bool) (lcs *api.Repository) {
	dep := *ptrDep
	if v, ok := cache[dep]; ok {
		return v
	}

	lcs = &api.Repository{}

	resp, err := http.Get(fmt.Sprintf("https://%v", dep))

	if err != nil {
		return
	}

	defer resp.Body.Close()

	data := new(MetaData)
	if err = m.Metabolize(resp.Body, data); err != nil {
		return
	}

	imports := strings.Split(data.Import, " ")
	if len(imports) != 3 {
		return
	}

	url := imports[2]
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 4 {
		return
	}
	if verbose {
		ptrDep = &urlParts[0]
	}
	lcs = &api.Repository{
		URL:    url,
		Host:   urlParts[2],
		Author: urlParts[3],
	}

	if len(urlParts) == 5 {
		lcs.Project = urlParts[4]
	}

	cache[dep] = lcs
	return
}
