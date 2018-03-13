package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/ribice/glice/api"
)

func TestGetCurrentFolder(t *testing.T) {
	s := strings.Split(getCurrentFolder(""), "src"+fs)[1]
	if s != "github.com"+fs+"ribice"+fs+"glice"+fs {
		t.Errorf("Current folder is not correct")
	}
}

func TestGetCurrentFolderWithPath(t *testing.T) {
	path := filepath.Join("github.com", "ribice", "glice", "testdata", "validate")
	s := strings.Split(getCurrentFolder(path), "src"+fs)[1]
	if s != path+fs {
		t.Errorf("Current folder is not correct")
	}
}
func TestGetFolders(t *testing.T) {
	cf, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	s := getFolders(cf+fs, "testdata")
	if s[0] != "." || s[1] != "api" {
		t.Errorf("There were missing or wrongly named folders")
	}
}

func TestGetFoldersWithPath(t *testing.T) {
	cf, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	s := getFolders(cf+fs+"testdata"+fs+"validate"+fs, "")
	if s[0] != "." || s[1] != "validators" {
		t.Errorf("There were missing or wrongly named folders")
	}
}
func TestGetLicenseWriteStd(t *testing.T) {
	ds := deps{
		deps: []dep{
			dep{
				name: "github.com/andygrunwald/go-jira",
				license: &api.License{
					URL:     "github.com/andygrunwald/go-jira",
					Author:  "andygrunwald",
					Project: "go-jira",
					Host:    "github.com",
				},
			},
		}}

	bd := "github.com/ribice/glice/testdata/validate/"
	bdl := len(bd) - 1
	ds.getDeps(bd, "."+fs, "Imports", bdl, true, false)
	if len(ds.deps) != 9 {
		t.Errorf("Incorrect number of dependencies")
	}
	if ds.deps[3].name != "github.com/markbates/going" {
		t.Errorf("Incorrect third dependency")
	}

	ds.getLicensesWriteStd("", nil, false, false, false)

	if ds.deps[0].license.Shortname != color.New(color.FgGreen).Sprintf("MIT") {
		t.Errorf("API did not return correct license.")
	}
}

func TestGetLicenseWriteFile(t *testing.T) {
	ds := deps{}

	cf, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := cf + fs + "testdata" + fs + "demo" + fs
	bd := strings.Split(path, "src"+fs)[1]
	bdl := len(bd) - 1
	ds.getDeps(bd, "."+fs, "Imports", bdl, true, false)

	ds.getLicensesWriteStd(path, nil, true, true, true)

	if _, err := os.Stat(path + "licenses" + fs); err != nil {
		if !os.IsNotExist(err) {
			t.Errorf("Folder licenses was not deleted")
		}
	}

}
