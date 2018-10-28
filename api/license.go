package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type licenseFormat struct {
	name  string
	color color.Attribute
}

var apis = map[string]string{
	"github.com": "https://api.github.com/repos/",
}

var licenseCol = map[string]licenseFormat{
	"other": {
		name:  "Other",
		color: color.FgBlue,
	},
	"mit": {
		name:  "MIT",
		color: color.FgGreen,
	},
	"lgpl-3.0": {
		name:  "LGPL-3.0",
		color: color.FgCyan,
	},
	"mpl-2.0": {
		name:  "MPL-2.0",
		color: color.FgHiBlue,
	},
	"agpl-3.0": {
		name:  "AGPL-3.0",
		color: color.FgHiCyan,
	},
	"unlicense": {
		name:  "Unlicense",
		color: color.FgHiRed,
	},
	"apache-2.0": {
		name:  "Apache-2.0",
		color: color.FgHiGreen,
	},
	"gpl-3.0": {
		name:  "GPL-3.0",
		color: color.FgHiMagenta,
	},
}

// License holds information about the license
type License struct {
	Shortname   string
	FullLicense string
	URL         string
	Exists      bool
	Host        string
	Author      string
	Project     string
}

// NewGitClient instantiates new GitClient
func NewGitClient(c context.Context, keys map[string]string) *GitClient {
	var tc *http.Client
	var ghLogged bool
	if v, ok := keys["github.com"]; ok {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: v},
		)
		tc = oauth2.NewClient(c, ts)
		ghLogged = true
	}
	ghClient := github.NewClient(tc)
	return &GitClient{GH: gHClient{
		cl: ghClient, logged: ghLogged,
	}}
}

// GitClient holds clients for interfering with Git provider APIs
type GitClient struct {
	GH gHClient
}

type gHClient struct {
	cl     *github.Client
	logged bool
}

// GetLicenses gets licenses for 3rd party dependencies
func (l *License) GetLicenses(c context.Context, gc *GitClient, repoStar, fileWrite bool) error {

	switch l.Host {
	case "github.com":
		if gc.GH.logged && repoStar {
			gc.GH.cl.Activity.Star(c, l.Author, l.Project)
		}
		rl, _, err := gc.GH.cl.Repositories.License(c, l.Author, l.Project)
		if err != nil {
			return err
		}
		name, clr := licenseCol[*rl.License.Key].name, licenseCol[*rl.License.Key].color
		if name == "" {
			name = *rl.License.Key
			clr = color.FgYellow
		}
		l.Shortname = color.New(clr).Sprintf(name)

		if fileWrite {
			l.writeToFile(rl.GetContent(), "licenses")
		}

	}
	return nil

}

func (l *License) writeToFile(s, folderName string) error {
	dec, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(folderName + string(filepath.Separator) + l.Author + "-" + l.Project + "-license" + ".MD")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
	return nil
}

// StarGlice stars Glice if user is logged in and thanks flag has been passed
func StarGlice(c context.Context, g *GitClient) {
	if g.GH.logged {
		g.GH.cl.Activity.Star(c, "ribice", "glice")
	}
}
