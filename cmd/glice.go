package main

import (
	"flag"
	"log"
	"os"

	"github.com/derekbassett/glice/v2"
)

func main() {
	var (
		fileWrite = flag.Bool("f", false, "Write all licenses to files")
		indirect  = flag.Bool("i", false, "Gets indirect modules as well")
		path      = flag.String("p", "", `Path of desired directory to be scanned with Glice (e.g. "github.com/ribice/glice")`)
		thx       = flag.Bool("t", false, "Stars dependent repos. Needs GITHUB_API_KEY env variable to work")
	)

	flag.Parse()

	if *path == "" {
		cf, err := os.Getwd()
		checkErr(err)
		*path = cf
	}

	cl, err := glice.NewClient(*path)
	checkErr(err)

	checkErr(cl.ParseDependencies(*indirect, *thx))

	cl.Print(os.Stdout)

	if *fileWrite {
		checkErr(cl.WriteLicensesToFile())
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
