package glice

import (
	"context"
	"testing"

	"github.com/fatih/color"
)

func TestGitHubAPINoKey(t *testing.T) {
	c := context.Background()
	l := &Repository{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}

	gc := newGitClient(c, map[string]string{}, false)
	err := gc.GetLicense(c, l)
	if err != nil {
		t.Error(err)
	}

	if l.Shortname != color.New(color.FgGreen).Sprintf("MIT") {
		t.Errorf("API did not return correct license or color.")
	}

}

func TestNonexistentLicense(t *testing.T) {

	c := context.Background()
	l := &Repository{
		URL:     "github.com/denysdovhan/wtfjs",
		Host:    "github.com",
		Author:  "denysdovhan",
		Project: "wtfjs",
	}

	gc := newGitClient(c, map[string]string{}, false)
	err := gc.GetLicense(c, l)
	if err != nil {
		t.Error(err)
	}

	if l.Shortname != color.New(color.FgYellow).Sprintf("wtfpl") {
		t.Errorf("API did not return correct license or color.")
	}

}

func TestGitHubAPIWithKey(t *testing.T) {

	c := context.Background()
	l := &Repository{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}

	v := map[string]string{
		"github.com": "apikey",
	}

	gc := newGitClient(c, v, false)
	err := gc.GetLicense(c, l)
	if err == nil {
		t.Error("expected bad credentials error")
	}

}

func TestGitHubAPIWithKeyAndThanks(t *testing.T) {

	c := context.Background()
	l := &Repository{
		URL:     "github.com/ribice/kiss",
		Host:    "github.com",
		Author:  "ribice",
		Project: "kiss",
	}

	v := map[string]string{
		"github.com": "apikey",
	}

	gc := newGitClient(c, v, true)

	err := gc.GetLicense(c, l)
	if err == nil {
		t.Error("expected bad credentials error")
	}

}
