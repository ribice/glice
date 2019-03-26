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
	name  string
	count int
	repo  *api.Repository
}

type deps struct {
	deps      []dep
	baseDir   string
	depth     string
	bdl       int
	incstdlib bool
	verbose   bool
	count     bool
	fw        bool
	cl        *api.GitClient
}

func main() {

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
		depth      = "Imports"
		apiKeys    = map[string]string{}
	)

	flag.Parse()

	fullPath := getCurrentFolder(*path)
	basedir := strings.Split(fullPath, "src"+fs)[1]
	bdl := len(basedir) - 1

	if *recursive {
		depth = "Deps"
	}

	ds := deps{
		baseDir:   basedir,
		bdl:       bdl,
		incstdlib: *incStdLib,
		verbose:   *verbose,
		count:     *count,
		fw:        *fileWrite,
		depth:     depth,
	}

	for _, v := range getFolders(fullPath, *ignoreDirs) {
		ds.getDeps(v)
	}

	if *ghkey != "" {
		apiKeys["github.com"] = *ghkey
	}
	c := context.Background()
	ds.cl = api.NewGitClient(c, apiKeys, *thx)
	ds.getLicenses(c)

	tw := tablewriter.NewWriter(os.Stdout)
	ds.writeStd(tw)
	tw.Render()
	checkErr(ds.writeLicensesToFile(fullPath))
}

func getCurrentFolder(path string) string {
	if path != "" {
		if !strings.HasSuffix(path, fs) {
			path += fs
		}
		return build.Default.GOPATH + fs + "src" + fs + path
	}
	cf, err := os.Getwd()
	checkErr(err)
	return cf + fs
}

func getFolders(fullPath, ignore string) []string {
	ign := strings.Split(ignore, ",")
	var folders []string
	checkErr(filepath.Walk(fullPath+".", func(path string, info os.FileInfo, err error) error {
		// Return only folders
		if info.IsDir() {
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
	}))
	return folders
}

func skipHidden(name string) bool {
	for _, v := range strings.Split(name, fs) {
		if strings.HasPrefix(v, ".") {
			return true
		}
	}
	return false
}

func (ds *deps) getDeps(dirname string) {
	// used for comparing dependency with current project minus the file separator
	if dirname == "."+fs {
		dirname = ""
	}
	args := "go list -f" + ` '{{ .` + ds.depth + ` }}' ` + ds.baseDir + dirname + ` | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' `
	if ds.incstdlib {
		args = "go list -f '{{  join ." + ds.depth + ` "\n"}}` + "' " + ds.baseDir + dirname
	}

	cmd := exec.Command("bash", "-c", args)
	stdout, err := cmd.StdoutPipe()
	checkErr(err)
	checkErr(cmd.Start())
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		if d := ds.exists(s.Text(), ds.verbose); d != nil {
			if len(d.name) >= ds.bdl && d.name[0:ds.bdl]+fs == ds.baseDir {
				continue
			}
			ds.deps = append(ds.deps, *d)
		}
	}
}

func (ds *deps) exists(s string, verbose bool) *dep {
	if strings.Contains(s, "vendor"+fs) {
		s = strings.Split(s, "vendor"+fs)[1]
	}

	l := getRepoURL(&s, verbose)

	for i, v := range ds.deps {
		if v.name == s {
			ds.deps[i].count++
			return nil
		}

		if v.repo != nil && l != nil && v.repo.URL == l.URL {
			l.Exists = true
		}
	}

	return &dep{name: s, repo: l}
}

func getRepoURL(s *string, verbose bool) *api.Repository {
	spl := strings.Split(*s, fs)
	switch spl[0] {
	case "github.com", "gitlab.com", "bitbucket.org":
		if len(spl) < 3 {
			return &api.Repository{}
		}
		if !verbose {
			*s = filepath.Join(spl[0], spl[1], spl[2])
		}
		return &api.Repository{URL: "https://" + spl[0] + "/" + spl[1] + "/" + spl[2], Host: spl[0], Author: spl[1], Project: spl[2]}

	case "gopkg.in":
		if len(spl) < 3 {
			return &api.Repository{}
		}
		if !verbose {
			*s = filepath.Join(spl[0], spl[1], spl[2])
		}
		return &api.Repository{URL: "https://github.com/" + spl[1] + "/" + strings.Split(spl[2], ".")[0], Host: "github.com", Author: spl[1], Project: strings.Split(spl[2], ".")[0]}
	}
	return getOtherRepo(s, verbose)
}

func (ds *deps) getLicenses(c context.Context) {
	for _, v := range ds.deps {
		v.repo.GetLicenses(c, ds.cl)
	}
}

func (ds *deps) writeStd(tw *tablewriter.Table) {
	keys := []string{"Dependency", "RepoURL", "License"}
	if ds.count {
		keys = append(keys, "Count")
	}
	tw.SetHeader(keys)
	for _, v := range ds.deps {
		vals := []string{v.name, color.BlueString(v.repo.URL), v.repo.Shortname}
		if ds.count {
			vals = append(vals, strconv.Itoa(v.count+1))
		}
		tw.Append(vals)
	}
}

func (ds *deps) writeLicensesToFile(path string) error {
	if !ds.fw {
		return nil
	}
	os.Mkdir("licenses", 0777)
	for _, v := range ds.deps {
		if v.repo.Text == "" {
			continue
		}

		dec, err := base64.StdEncoding.DecodeString(v.repo.Text)
		if err != nil {
			return err
		}
		f, err := os.Create(path + fs + v.repo.Author + "-" + v.repo.Project + "-license.MD")
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
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
