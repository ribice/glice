package main

import (
	"testing"
)

func TestGetIndirectRepo(t *testing.T) {
	indirectRepo := getIndirectRepo("golang.org/x/net/context/ctxhttp")

	if indirectRepo.URL != "https://go.googlesource.com/net" {
		t.Errorf("Wrong URL: %s", indirectRepo.URL)
	}
}
