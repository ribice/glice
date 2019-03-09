package main

import (
	"testing"
)

func TestGetOtherRepo(t *testing.T) {
	str := "golang.org/x/net/context/ctxhttp"
	r := getOtherRepo(&str, false)

	if r.URL != "https://go.googlesource.com/net" {
		t.Errorf("Wrong URL: %s", r.URL)
	}
}
