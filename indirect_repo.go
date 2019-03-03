package main

// Resolve indirect repos as described here:
// https://golang.org/cmd/go/#hdr-Remote_import_paths

import (
	"errors"
	"net/http"
	"strings"

	m "github.com/keighl/metabolize"
)

type MetaData struct {
	Import string `meta:"go-import"`
	Source string `meta:"go-source"`
}

type IndirectRepo struct {
	Dep     string
	Repo    string
	URL     string
	Author  string
	Project string
}

func getIndirectRepo(dep string) (indirectRepo IndirectRepo, err error) {

	path := "https://" + dep

	resp, err := http.Get(path)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	data := new(MetaData)
	err = m.Metabolize(resp.Body, data)
	if err != nil {
		return
	}

	indirectRepo, err = headerToIndirectRepo(data.Import)
	return
}

func headerToIndirectRepo(importHeader string) (indirectRepo IndirectRepo, err error) {
	if importHeader == "" {
		err = errors.New("header is empty")
		return
	}

	frags := strings.Split(importHeader, " ")
	if len(frags) != 3 {
		err = errors.New("header has wrong format: " + importHeader)
		return
	}

	URL := frags[2]
	repo, author, project, err := parseURL(URL)
	if err != nil {
		return
	}

	indirectRepo = IndirectRepo{
		Dep:     frags[0],
		Repo:    repo,
		URL:     frags[2],
		Author:  author,
		Project: project,
	}

	return
}

func parseURL(URL string) (repo string, author string, project string, err error) {
	frags := strings.Split(URL, "/")
	len := len(frags)

	if len < 4 {
		err = errors.New("invalid URL: " + URL)
	} else {
		repo = frags[2]
		author = frags[3]
		if len == 5 {
			project = frags[4]
		}
	}

	return
}
