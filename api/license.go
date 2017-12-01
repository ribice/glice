package api

import (
	"context"
	"encoding/base64"
	"fmt"
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
	"other": licenseFormat{
		name:  "Other",
		color: color.FgBlue,
	},
	"mit": licenseFormat{
		name:  "MIT",
		color: color.FgGreen,
	},
	"lgpl-3.0": licenseFormat{
		name:  "LGPL-3.0",
		color: color.FgCyan,
	},
	"mpl-2.0": licenseFormat{
		name:  "MPL-2.0",
		color: color.FgHiBlue,
	},
	"agpl-3.0": licenseFormat{
		name:  "AGPL-3.0",
		color: color.FgHiCyan,
	},
	"unlicense": licenseFormat{
		name:  "Unlicense",
		color: color.FgHiRed,
	},
	"apache-2.0": licenseFormat{
		name:  "Apache-2.0",
		color: color.FgHiGreen,
	},
	"gpl-3.0": licenseFormat{
		name:  "GPL-3.0",
		color: color.FgHiMagenta,
	},
}

// License holds information about the license
type License struct {
	Shortname   string
	fullLicense string
	URL         string
	Exists      bool
	Host        string
	Author      string
	Project     string
}

// GetLicenses gets licenses for 3rd party dependencies
func (l *License) GetLicenses(c context.Context, keys map[string]string, fileWrite bool) error {

	var tc *http.Client

	switch l.Host {
	case "github.com":
		if v, ok := keys["github.com"]; ok {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: v},
			)
			tc = oauth2.NewClient(c, ts)
		}
		client := github.NewClient(tc)

		rl, _, err := client.Repositories.License(c, l.Author, l.Project)
		if err != nil {
			return fmt.Errorf("bad credentials")
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
