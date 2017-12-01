package main

import (
	"fmt"
	"testing"

	"github.com/ribice/glice/api"

	"github.com/fatih/color"
)

func TestGetCurrentFolder(t *testing.T) {
	s := getCurrentFolder()
	if s != "github.com/ribice/glice/" {
		t.Errorf("Current folder is not correct")
	}
}
func TestGetFolders(t *testing.T) {
	s := getFolders("testdata")
	if s[0] != "."+fs || s[1] != "api" {
		t.Errorf("There were missing or wrongly named folders")
	}

}

func TestGetLicenseWriteStd(t *testing.T) {
	ds := deps{}
	d := dep{
		name: "github.com/andygrunwald/go-jira",
		license: &api.License{
			URL:     "github.com/andygrunwald/go-jira",
			Author:  "andygrunwald",
			Project: "go-jira",
			Host:    "github.com",
		},
	}

	ds.deps = append(ds.deps, d)
	bd := "github.com/ribice/glice/testdata/validate/"
	bdl := len(bd) - 1
	ds.getDeps(bd, "."+fs, "Imports", bdl, true, false)
	if len(ds.deps) != 9 {
		t.Errorf("Incorrect number of dependencies")
	}
	if ds.deps[3].name != "github.com/markbates/going" {
		t.Errorf("Incorrect third dependency")
	}

	ds.getLicensesWriteStd(nil, false)
	fmt.Println(ds.deps[0].license)
	if ds.deps[0].license.Shortname != color.New(color.FgGreen).Sprintf("MIT") {
		t.Errorf("API did not return correct license.")
	}

}
