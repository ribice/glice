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

func getOtherRepo(ptrDep *string, verbose bool) *api.Repository {
	dep := *ptrDep
	if v, ok := cache[dep]; ok {
		return v
	}
	resp, err := http.Get(fmt.Sprintf("https://%v", dep))

	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	data := new(MetaData)
	if err = m.Metabolize(resp.Body, data); err != nil {
		return nil
	}

	imports := strings.Split(data.Import, " ")
	if len(imports) != 3 {
		return nil
	}

	url := imports[2]
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 4 {
		return nil
	}
	if verbose {
		ptrDep = &urlParts[0]
	}
	lcs := &api.Repository{
		URL:    url,
		Host:   urlParts[2],
		Author: urlParts[3],
	}

	if len(urlParts) == 5 {
		lcs.Project = urlParts[4]
	}

	cache[dep] = lcs
	return lcs
}
