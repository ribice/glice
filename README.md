# glice

[![Build Status](https://img.shields.io/github/workflow/status/ribice/glice/CI?style=flat-square)](https://github.com/ribice/glice/actions?query=workflow%3ACI)
[![Coverage Status](https://coveralls.io/repos/github/ribice/glice/badge.svg?branch=master)](https://coveralls.io/github/ribice/glice?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/glice)](https://goreportcard.com/report/github.com/ribice/glice)

Golang license and dependency checker. Prints list of all dependencies, their URL, license and saves all the license files in /licenses.

## Introduction

glice analyzes the go.mod file of your project and prints it in a tabular format - name, URL, and license short-name (MIT, GPL...).

## Installation

Download and install glice by executing:

```bash
    go get github.com/ribice/glice/v2
    go install github.com/ribice/glice/v2/cmd/glice.go
```

To update:

```bash
    go get -u github.com/ribice/glice/v2
```

## Usage

To run glice, navigate to a folder with go.mod and execute:

```bash
    glice
```

Alternatively, you can provide path which you want to be scanned with -p flag:

```bash
    glice -p "github.com/ribice/glice"
```

By default glice:

- Prints only to stdout

- Gets dependencies from  go.mod

- Fetch licenses for dependencies hosted on GitHub
  
- Is limited to 60 API calls on GitHub (up to 60 dependencies from github.com). API key can be provided by setting `GITHUB_API_KEY` environment variable.

All flags are optional. Glice supports the following flags:

```
- f [boolean, fileWrite] // Writes all licenses to /licenses dir
- i [boolean, indirect] // Parses indirect dependencies as well
- p [string - path] // Path to be scanned in form of github.com/author/repo
- t [boolean - thanks] // if GitHub API key is provided, setting this flag will star all GitHub repos from dependency. __In order to do this, API key must have access to public_repo__

```

Don't forget `-help` flag for detailed usage information.

## Using glice inside as a library

As of v2.0.0 glice can be used as a library and provides few functions/methods that return list of dependencies in structured format and printing to io.Writer.

## Sample output

Executing glice -c on github.com/ribice/glice prints (with additional colors for links and licenses):

```
+-----------------------------------+-------------------------------------------+--------------+
|            DEPENDENCY             |                  REPOURL                  |   LICENSE    |
+-----------------------------------+-------------------------------------------+--------------+
| github.com/fatih/color            | https://github.com/fatih/color            | MIT          |
| github.com/google/go-github       | https://github.com/google/go-github       | bsd-3-clause |
| github.com/keighl/metabolize      | https://github.com/keighl/metabolize      | Other        |
| github.com/olekukonko/tablewriter | https://github.com/olekukonko/tablewriter | MIT          |
| golang.org/x/mod                  | https://go.googlesource.com/mod           |              |
| golang.org/x/oauth2               | https://go.googlesource.com/oauth2        |              |
+-----------------------------------+-------------------------------------------+--------------+
```

## License

glice is licensed under the MIT license. Check the [LICENSE](LICENSE.md) file for details.

## Author

[Emir Ribic](https://ribice.ba)