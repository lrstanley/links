<p align="center">links -- Simple, fast link shortener</p>
<p align="center">
  <a href="https://travis-ci.org/lrstanley/links"><img src="https://travis-ci.org/lrstanley/links.svg?branch=master" alt="Build Status"></a>
  <a href="https://byteirc.org/channel/%23%2Fdev%2Fnull"><img src="https://img.shields.io/badge/ByteIRC-%23%2Fdev%2Fnull-blue.svg" alt="IRC Chat"></a>
</p>

## Table of Contents
- [Installation](#installation)
  - [Ubuntu/Debian](#ubuntudebian)
  - [CentOS/Redhat](#centosredhat)
  - [Manual Install](#manual-install)
  - [Build from source](#build-from-source)
- [Usage](#usage)
  - [Example](#example)
- [Using as a library](#using-as-a-library)
  - [Example](#example-1)
- [Migrating](#migrating)
- [API](#api)
    - [Password protection](#password-protection)
- [Contributing](#contributing)
- [License](#license)

## Installation

Check out the [releases](https://github.com/lrstanley/links/releases)
page for prebuilt versions. Links should work on ubuntu/debian,
centos/redhat/fedora, etc. Below are example commands of how you would install
the utility.

### Ubuntu/Debian

```bash
$ wget https://liam.sh/ghr/links_[[tag]]_[[os]]_[[arch]].deb
$ dpkg -i links_[[tag]]_[[os]]_[[arch]].deb
$ links --help
```

### CentOS/Redhat

```bash
$ yum localinstall https://liam.sh/ghr/links_[[tag]]_[[os]]_[[arch]].rpm
$ links --help
```

Some older CentOS versions may require (if you get `Cannot open: <url>. Skipping.`):

```console
$ wget https://liam.sh/ghr/links_[[tag]]_[[os]]_[[arch]].rpm
$ yum localinstall links_[[tag]]_[[os]]_[[arch]].rpm
```

### Manual Install

```bash
$ wget https://liam.sh/ghr/links_[[tag]]_[[os]]_[[arch]].tar.gz
$ tar -C /usr/bin/ -xzvf links_[[tag]]_[[os]]_[[arch]].tar.gz links
$ chmod +x /usr/bin/links
$ links --help
```

### Source

Note that you must have [Go](https://golang.org/doc/install) installed and
a fully working `$GOPATH` setup.

    $ go get -d -u github.com/lrstanley/links
    $ cd $GOPATH/src/github.com/lrstanley/links
    $ make
    $ ./links --help

## Usage

```
$ links --help
Usage:
  links [OPTIONS]

Application Options:
  -s, --site-name=    site url, used for url generation (default: https://links.wtf)
  -q, --quiet         don't log to stdout
  -b, --http=         ip:port pair to bind to (default: :8080)
  -p, --behind-proxy  if X-Forwarded-For headers should be trusted
  -d, --db=           path to database file (default: store.db)
      --key-length=   default length of key (uuid) for generated urls (default: 4)
  -e, --export-file=  file to export db to (default: links.export)
      --export-json   export db to json elements
      --migrate       begin migration from Links running MySQL
      --migrate-info= connection url used to connect to the old mysql instance (default: user:passwd@tcp(localhost:3306)/links_db)
  -v, --version       display the version of Links and exit

TLS Options:
      --tls.enable    run tls server rather than standard http
  -c, --tls.cert=     path to ssl cert file
  -k, --tls.key=      path to ssl key file

Help Options:
  -h, --help          Show this help message
```

### Example

```
$ links -s "http://your-domain.com" -b "0.0.0.0:80" -d links.db
```

## Using as a library

Links also has a Go client library which you can use, which adds a simple
wrapper around an http call, to make shortening links simpler. Download it
using the following `go get` command:

```
$ go get -u github.com/lrstanley/links/client
```

View the documentation [here](https://godoc.org/github.com/lrstanley/links/client)

### Example

```go
package main

import (
	"fmt"
	"log"

	links "github.com/lrstanley/links/client"
)

func main() {
	uri, err := links.Shorten("https://your-long-link.com/longer/link", "", nil)
	if err != nil {
		log.Fatalf("unable to shorten link: %s", err)
	}

	fmt.Printf("shortened: %s\n", uri.String())
}
```

## Migrating

If you were using Links before it was rewritten in Go (and was using
MySQL), then this is how you can migrate the database:

```
$ links --migrate --migrate-info "user:passwd@tcp(your-server:3306)/links_db"
```

Depending on the amount of links within the old database, this may take a while
to complete. I would recommend backing things up prior just in case as well.

## API

Shortening a link is quite easy. simply send a `POST` request to `https://links.wtf/add`,
which will return JSON-safe information as shown below:

```
$ curl --data "url=http://google.com" https://links.wtf/add
{"url": "https://links.wtf/27f4", "success": true}
```

#### Password protection

You can also password protect a link, simply by adding a `password` variable to the payload:

```
$ curl --data 'url=https://google.com/example&encrypt=y0urp4$$w0rd' https://links.wtf/add
{"url": "https://links.wtf/abc123", "success": true}
```

## Contributing

Please review the [CONTRIBUTING](CONTRIBUTING.md) doc for submitting issues/a guide
on submitting pull requests and helping out.

## License

```
LICENSE: The MIT License (MIT)
Copyright (c) 2016 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
