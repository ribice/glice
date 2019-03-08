package api

import (
	"context"
	"net/http"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type licenseFormat struct {
	name  string
	color color.Attribute
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

// Repository holds information about the repository
type Repository struct {
	Shortname   string
	fullLicense string
	URL         string
	Exists      bool
	Host        string
	Author      string
	Project     string
	Text        string
}

// NewGitClient instantiates new GitClient
func NewGitClient(c context.Context, keys map[string]string, star bool) *GitClient {
	var tc *http.Client
	var ghLogged bool
	if v, ok := keys["github.com"]; ok {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: v},
		)
		tc = oauth2.NewClient(c, ts)
		ghLogged = true
	}
	return &GitClient{GH: gHClient{
		cl:     github.NewClient(tc),
		logged: ghLogged,
	}, star: star}
}

// GitClient holds clients for interfering with Git provider APIs
type GitClient struct {
	GH   gHClient
	star bool
}

type gHClient struct {
	cl     *github.Client
	logged bool
}

// GetLicenses gets licenses for 3rd party dependencies
func (r *Repository) GetLicenses(c context.Context, gc *GitClient) {
	switch r.Host {
	case "github.com":
		rl, _, err := gc.GH.cl.Repositories.License(c, r.Author, r.Project)
		if err == nil {
			name, clr := licenseCol[*rl.License.Key].name, licenseCol[*rl.License.Key].color
			if name == "" {
				name = *rl.License.Key
				clr = color.FgYellow
			}
			r.Shortname = color.New(clr).Sprintf(name)
			r.Text = rl.GetContent()
		}
		if gc.star {
			gc.GH.cl.Activity.Star(c, r.Author, r.Project)
		}
	}
}
