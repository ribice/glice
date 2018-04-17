package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/ribice/glice/api"
)

const (
	nl = "\n"
	fs = string(filepath.Separator)
)

type dep struct {
	name    string
	count   int
	license *api.License
}

type deps struct {
	deps []dep
}

func main() {

	var ds deps

	var (
		verbose    = flag.Bool("v", false, "Include detailed imports (github.com/author/repo/net/http, github.com/author/repo/net/middleware ... instead of only github.com/author/repo")
		incStdLib  = flag.Bool("s", false, "Include standard library dependencies")
		recursive  = flag.Bool("r", false, "Gets single level recursive dependencies")
		fileWrite  = flag.Bool("f", false, "Writes all licenses to files")
		ignoreDirs = flag.String("i", "", "Comma separated list of folders that should be ignored")
		ghkey      = flag.String("gh", "", "GitHub API key used for increasing the GitHub's API rate limit from 60req/h to 5000req/h")
		path       = flag.String("p", "", `Path of desired directory to be scanned with Glice (e.g. "github.com/ribice/glice/")`)
		thx        = flag.Bool("t", false, "Stars dependent repos.")
		count      = flag.Bool("c", false, "Include usage count in exported result")
		indirect   = flag.Bool("in", false, "Resolve indirect repos (find the repo location through an html meta header)")
		depth      = "Imports"
		apiKeys    = map[string]string{}
	)

	flag.Parse()

	// Gets current folder in $GOPATH

	fullPath := getCurrentFolder(*path)
	basedir := strings.Split(fullPath, "src"+fs)[1]
	bdl := len(basedir) - 1

	if *recursive {
		depth = "Deps"
	}

	if *ghkey != "" {
		apiKeys["github.com"] = *ghkey
	}

	for _, v := range getFolders(fullPath, *ignoreDirs) {
		// implement concurrency here
		ds.getDeps(basedir, v, depth, bdl, *incStdLib, *verbose, *indirect)
	}
	c := context.Background()
	gh := api.NewGitClient(c, apiKeys)
	ds.getLicenses(c, gh)
	if *thx {
		ds.starRepos(c, gh)
	}
	tw := tablewriter.NewWriter(os.Stdout)
	if !*count {
		ds.writeStd(tw)
	} else {
		ds.writeStdCount(tw)
	}
	tw.Render()
	if *fileWrite {
		if err := ds.writeLicensesToFile(fullPath); err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}
	}
}

func getCurrentFolder(path string) string {
	if path != "" {
		if !strings.HasSuffix(path, fs) {
			path += fs
		}
		return build.Default.GOPATH + fs + "src" + fs + path
	}
	cf, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cf + fs
}

func getFolders(fullPath, ignore string) []string {
	ign := strings.Split(ignore, ",")
	var folders []string
	err := filepath.Walk(fullPath+".", func(path string, info os.FileInfo, err error) error {
		// Return only folders
		if info.IsDir() {
			//name := strings.Split(info.Name(), "src"+fs)[1]
			// Skip if folder name is vendor, is hidden (starting with dot, but ignore dot only)
			if (info.Name() == "vendor" || skipHidden(info.Name())) && info.Name() != "." {
				return filepath.SkipDir
			}
			for _, v := range ign {
				if info.Name() == v {
					return filepath.SkipDir
				}
			}

			folders = append(folders, strings.Split(path, fullPath)[1])

		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return folders
}

func skipHidden(name string) bool {
	split := strings.Split(name, fs)
	for _, v := range split {
		if strings.HasPrefix(v, ".") == true {
			return true
		}
	}
	return false
}

func (ds *deps) getDeps(basedir, dirname, depth string, bdl int, incStdLib, verbose bool, indirect bool) {

	// used for comparing dependency with current project minus the file separator
	if dirname == "."+fs {
		dirname = ""
	}
	args := "go list -f" + ` '{{ .` + depth + ` }}' ` + basedir + dirname + ` | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' `
	if incStdLib {
		args = "go list -f '{{  join ." + depth + ` "\n"}}` + "' " + basedir + dirname
	}

	cmd := exec.Command("bash", "-c", args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		if d := ds.exists(s.Text(), verbose, indirect); d != nil {
			if len(d.name) >= bdl && d.name[0:bdl]+fs == basedir {
				continue
			}
			ds.deps = append(ds.deps, *d)
		}
	}
}

func (ds *deps) exists(s string, verbose, indirect bool) *dep {
	if strings.Contains(s, "vendor"+fs) {
		s = strings.Split(s, "vendor"+fs)[1]
	}

	l := getRepoURL(&s, verbose, indirect)

	for i, v := range ds.deps {
		if v.name == s {
			ds.deps[i].count++
			return nil
		}

		if v.license != nil && l != nil && v.license.URL == l.URL {
			l.Exists = true
		}
	}

	return &dep{name: s, license: l}
}

func getRepoURL(s *string, verbose bool, indirect bool) *api.License {
	spl := strings.Split(*s, fs)
	switch spl[0] {
	case "github.com", "gitlab.com", "bitbucket.org":
		if !verbose && len(spl) >= 3 {
			*s = filepath.Join(spl[0], spl[1], spl[2])
		}
		if len(spl) >= 3 {
			return &api.License{URL: "https://" + spl[0] + "/" + spl[1] + "/" + spl[2], Host: spl[0], Author: spl[1], Project: spl[2]}
		}
		return nil
	case "gopkg.in":
		if !verbose && len(spl) >= 3 {
			*s = filepath.Join(spl[0], spl[1], spl[2])
		}
		if len(spl) >= 3 {
			return &api.License{URL: "https://" + "github.com/" + spl[1] + "/" + strings.Split(spl[2], ".")[0], Host: "github.com", Author: spl[1], Project: strings.Split(spl[2], ".")[0]}
		}
		return nil
	default:
		if indirect {
			repo := getIndirectRepo(*s)
			if repo.Found {
				if !verbose {
					*s = repo.Dep
				}
				return &api.License{URL: repo.URL, Host: repo.Repo, Author: repo.Author, Project: repo.Project}
			}
		}
		return nil
	}
}

func (ds *deps) getLicenses(c context.Context, gh *api.GitClient) {
	for _, v := range ds.deps {
		if v.license != nil {
			v.license.GetLicenses(c, gh)
		}
	}
}

func (ds *deps) starRepos(c context.Context, gh *api.GitClient) {
	for _, v := range ds.deps {
		if v.license != nil {
			v.license.Star(c, gh)
		}
	}
	api.StarGlice(c, gh)
}

func (ds *deps) writeStd(tw *tablewriter.Table) {
	tw.SetHeader([]string{"Dependency", "RepoURL", "License"})
	for _, v := range ds.deps {
		str := []string{v.name, strconv.Itoa(v.count + 1)}
		if v.license != nil {
			str = append(str, color.BlueString(v.license.URL), v.license.Shortname)
		} else {
			str = append(str, "", "")
		}
		tw.Append(str)
	}
}

func (ds *deps) writeStdCount(tw *tablewriter.Table) {
	tw.SetHeader([]string{"Dependency", "Count", "RepoURL", "License"})
	for _, v := range ds.deps {
		str := []string{v.name}
		if v.license != nil {
			str = append(str, color.BlueString(v.license.URL), v.license.Shortname)
		} else {
			str = append(str, "", "")
		}
		tw.Append(str)
	}
}

func (ds *deps) writeLicensesToFile(path string) error {
	os.Mkdir("licenses", 0777)
	for _, v := range ds.deps {
		if v.license != nil && v.license.Text != "" {
			dec, err := base64.StdEncoding.DecodeString(v.license.Text)
			if err != nil {
				return err
			}
			f, err := os.Create(path + fs + v.license.Author + "-" + v.license.Project + "-license" + ".MD")
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := f.Write(dec); err != nil {
				return err
			}
			if err := f.Sync(); err != nil {
				return err
			}
		}
	}
	return nil
}
