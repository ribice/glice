package glice

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	m "github.com/keighl/metabolize"
	"github.com/olekukonko/tablewriter"

	"github.com/derekbassett/glice/v2/mod"
)

var (
	// ErrNoGoMod is returned when path doesn't contain go.mod file
	ErrNoGoMod = errors.New("no go.mod file present")

	// ErrNoAPIKey is returned when thanks flag is enabled without providing GITHUB_API_KEY env variable
	ErrNoAPIKey = errors.New("cannot use thanks feature without github api key")
)

type Client struct {
	dependencies []*Repository
	path         string
}

func NewClient(path string) (*Client, error) {
	if !mod.Exists(path) {
		return nil, ErrNoGoMod
	}
	return &Client{path: path}, nil
}

func (c *Client) ParseDependencies(includeIndirect, thanks bool) error {
	githubAPIKey := os.Getenv("GITHUB_API_KEY")
	if thanks && githubAPIKey == "" {
		return ErrNoAPIKey
	}
	repos, err := ListRepositories(c.path, includeIndirect)
	if err != nil {
		return err
	}

	ctx := context.Background()
	gitCl := newGitClient(ctx, map[string]string{"github.com": githubAPIKey}, thanks)
	for _, r := range repos {
		err = gitCl.GetLicense(ctx, r)
		if err != nil {
			log.Println(err)
		}
	}
	c.dependencies = repos
	return nil
}

func (c *Client) Print(output io.Writer) {
	if len(c.dependencies) < 1 {
		return
	}
	tw := tablewriter.NewWriter(output)
	tw.SetHeader([]string{"Dependency", "RepoURL", "License"})
	for _, d := range c.dependencies {
		tw.Append([]string{d.Name, color.BlueString(d.URL), d.Shortname})
	}
	tw.Render()
}

func Print(path string, indirect bool, writeTo io.Writer) error {
	c, err := NewClient(path)
	if err != nil {
		return err
	}

	err = c.ParseDependencies(indirect, false)
	if err != nil {
		return err
	}

	c.Print(writeTo)
	return nil
}

func ListRepositories(path string, withIndirect bool) ([]*Repository, error) {
	modules, err := mod.Parse(path, withIndirect)
	if err != nil {
		return nil, err
	}

	var repos []*Repository

	for _, mods := range modules {
		repos = append(repos, getRepository(mods))
	}

	return repos, nil

}

func getRepository(s string) *Repository {
	spl := strings.Split(s, "/")
	switch spl[0] {
	case "github.com", "gitlab.com", "bitbucket.org":
		if len(spl) < 3 {
			return &Repository{Name: s}
		}
		return &Repository{URL: "https://" + spl[0] + "/" + spl[1] + "/" + spl[2], Host: spl[0], Author: spl[1], Project: spl[2], Name: s}

	case "gopkg.in":
		if len(spl) < 3 {
			return &Repository{Name: s}
		}
		return &Repository{URL: "https://github.com/" + spl[1] + "/" + strings.Split(spl[2], ".")[0], Host: "github.com", Author: spl[1], Project: strings.Split(spl[2], ".")[0], Name: s}
	}
	return getOtherRepo(s)
}

type metaData struct {
	Import string `meta:"go-import"`
	Source string `meta:"go-source"`
}

var cache = map[string]*Repository{}

// Resolve indirect repos as described here:
// https://golang.org/cmd/go/#hdr-Remote_import_paths
func getOtherRepo(name string) *Repository {
	if v, ok := cache[name]; ok {
		return v
	}

	lcs := &Repository{Name: name}

	resp, err := http.Get(fmt.Sprintf("https://%s", name))
	if err != nil {
		return lcs
	}

	defer resp.Body.Close()

	data := new(metaData)
	if err = m.Metabolize(resp.Body, data); err != nil {
		return lcs
	}

	imports := strings.Split(data.Import, " ")
	if len(imports) != 3 {
		return lcs
	}

	url := imports[2]
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 4 {
		return lcs
	}

	lcs.URL = strings.TrimSuffix(url, ".git")
	lcs.Host = urlParts[2]
	lcs.Author = urlParts[3]

	if len(urlParts) == 5 {
		lcs.Project = strings.TrimSuffix(urlParts[4], ".git")
	}

	cache[name] = lcs
	return lcs
}

func (c *Client) WriteLicensesToFile() error {
	if len(c.dependencies) < 1 {
		return nil
	}
	os.Mkdir("licenses", 0777)

	for _, d := range c.dependencies {
		if d.Text == "" {
			continue
		}

		dec, err := base64.StdEncoding.DecodeString(d.Text)
		if err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(c.path, "licenses", d.Author+"-"+d.Project+"-license.MD"))
		if err != nil {
			return err
		}

		if _, err := f.Write(dec); err != nil {
			return err
		}
		if err := f.Sync(); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}
