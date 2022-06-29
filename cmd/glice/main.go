package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ribice/glice/v2"
)

func main() {
	var (
		fileWrite = flag.Bool("f", false, "Write all licenses to files")
		indirect  = flag.Bool("i", false, "Gets indirect modules as well")
		path      = flag.String("p", "", `Path of desired directory to be scanned with Glice (e.g. "github.com/ribice/glice/v2")`)
		thx       = flag.Bool("t", false, "Stars dependent repos. Needs GITHUB_API_KEY env variable to work")
		verbose   = flag.Bool("v", false, "Adds verbose logging")
		format    = flag.String("fmt", "table", "Output format [table | json | csv]")
		output    = flag.String("o", "stdout", "Output location [stdout | file]")
		extension = map[string]string{
			"table": "txt",
			"json":  "json",
			"csv":   "csv",
		}
	)

	flag.Parse()

	if *path == "" {
		cf, err := os.Getwd()
		checkErr(err)
		*path = cf
	}

	if !*verbose {
		log.SetOutput(io.Discard)
		log.SetFlags(0)
	}

	cl, err := glice.NewClient(*path, *format, *output)
	checkErr(err)

	checkErr(cl.ParseDependencies(*indirect, *thx))

	switch *output {
	case "stdout":
		cl.Print(os.Stdout)
	case "file":
		fileName := fmt.Sprintf("dependencies.%s", extension[*format])
		f, err := os.Create(fileName)
		checkErr(err)
		cl.Print(f)
		f.Close()
	}

	if *fileWrite {
		checkErr(cl.WriteLicensesToFile())
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
