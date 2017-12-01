# glice

[![Build Status](https://travis-ci.org/ribice/glice.svg?branch=master)](https://travis-ci.org/ribice/glice)
[![Coverage Status](https://coveralls.io/repos/github/ribice/glice/badge.svg?branch=master)](https://coveralls.io/github/ribice/glice?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribice/glice)](https://goreportcard.com/report/github.com/ribice/glice)

Golang license and dependency checker. Prints list of all dependencies (both from std and 3rd party), number of times used, their license and saves all the license files in /licenses.

## Introduction

glice analyzes the source code of your project and gets the list of all dependencies (by default only third-party ones, standard library can be enabled with flag) and prints in a tabular format their name, count (in different packages), URL (for third-party ones) and license short-name (MIT, GPL..).

## Installation

Download and install glice by executing:

```bash
    go get github.com/ribice/glice
    go install github.com/ribice/glice
```

To update:

```bash
    go get -u github.com/ribice/glice
```

## Usage

To run glice, navigate to a folder in gopath and execute:

```bash
    glice
```

By default glice:

- Prints only to stdout

- Shows only third-party dependencies, grouped under project name (e.g. github.com/author/repo/http and github.com/author/repo/http/middleware are counted as 2 github.com/author/repo, and that's being displayed)

- Shows only dependencies of project where you're currently located (does not go recursively through dependencies)

- Is limited to 60 API calls on GitHub (up to 60 dependencies from github.com)

- Ignores files in vendor folder

All flags are optional. Glice supports the following flags:

```
- s [boolean - include standard library]// return standard library dependencies as well as third-party ones
- v [boolean - verbose] // by default, github.com/author/repo/http and github.com/author/repo/http/middleware are counted as 2 times github.com/author/repo, and that's being displayed. With verbose, they are counted and shown separately.
- r [boolean - recursive] // returns dependencies of dependencies (single level only)
- f [boolean - fileWrite] // writes all license files inside /licenses folder
- i [string - ignoreFolders] // list of comma-separated folders that should be ignored
- gh [string - githubAPIKey] // used to increase GitHub API's rate limit from 60req/h to 5000req/h
```

Don't forget `-help` flag for detailed usage information.

## Sample output

Executing glice on github.com/qiangxue/golang-restful-starter-kit prints (with additional colors for links and licenses):

```
+------------------------------------+-------+--------------------------------------------+---------+
|             DEPENDENCY             | COUNT |                  REPOURL                   | LICENSE |
+------------------------------------+-------+--------------------------------------------+---------+
| fmt                                |     5 |
| github.com/Sirupsen/logrus         |     2 | https://github.com/Sirupsen/logrus         | MIT     |
| github.com/go-ozzo/ozzo-dbx        |     3 | https://github.com/go-ozzo/ozzo-dbx        | MIT     |
| github.com/go-ozzo/ozzo-routing    |     9 | https://github.com/go-ozzo/ozzo-routing    | MIT     |
| github.com/lib/pq                  |     2 | https://github.com/lib/pq                  | MIT     |
| net/http                           |     3 |
| github.com/dgrijalva/jwt-go        |     1 | https://github.com/dgrijalva/jwt-go        | MIT     |
| strconv                            |     1 |
| time                               |     2 |
| database/sql                       |     1 |
| github.com/go-ozzo/ozzo-validation |     3 | https://github.com/go-ozzo/ozzo-validation | MIT     |
| github.com/spf13/viper             |     1 | https://github.com/spf13/viper             | MIT     |
| gopkg.in/yaml.v2                   |     1 |
| io/ioutil                          |     2 |
| sort                               |     1 |
| strings                            |     3 |
| os                                 |     1 |
+------------------------------------+-------+--------------------------------------------+---------+
```

## To-Do

- [ ] Improve tests and code coverage
- [ ] Add flag to disable dependencies from test files
- [ ] Implement license checking for projects hosted on GitLab.com
- [ ] Implement license checking for projects hosted on Bitbucket.org
- [ ] Implement license checking for projects hosted on other third party sites (e.g. gitea.io, gogs.io)
- [ ] Add ability to send path of project via flag
- [ ] Remove dependency on go-github
- [ ] Implement concurrency
- [ ] Add naming option when saving licenses to files (currently author-repo-license.MD)

## License

glice is licensed under the MIT license. Check the [LICENSE](LICENSE.md) file for details.

## Author

[Emir Ribic](https://ribice.ba)