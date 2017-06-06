<p align="center">links.ml -- Simple, fast link shortener</p>
<p align="center">
  <a href="https://travis-ci.org/lrstanley/links.ml"><img src="https://travis-ci.org/lrstanley/links.ml.svg?branch=master" alt="Build Status"></a>
  <a href="http://byteirc.org/channel/L"><img src="https://img.shields.io/badge/ByteIRC-%23L-blue.svg" alt="IRC Chat"></a>
</p>

## Installing

### Prebuilt archives/binaries

Feel free to head over to the [releases](https://github.com/lrstanley/links.ml/releases)
page to get the latest available version of links.ml

### Build from source

    $ go get -u github.com/lrstanley/links.ml
    $ $GOPATH/bin/link.ml --help

## Usage

```
$ links --help
Usage:
  links [OPTIONS]

Application Options:
  -s, --site-name=    site url, used for url generation (default: https://links.ml)
  -q, --quiet         don't log to stdout
  -b, --http=         ip:port pair to bind to (default: :8080)
  -p, --behind-proxy  if X-Forwarded-For headers should be trusted
  -d, --db=           path to database file (default: store.db)
      --key-length=   default length of key (uuid) for generated urls (default: 4)
  -e, --export-file=  file to export db to (default: links.export)
      --export-json   export db to json elements
      --migrate       begin migration from links.ml running MySQL
      --migrate-info= connection url used to connect to the old mysql instance (default: user:passwd@tcp(localhost:3306)/links_db)
  -v, --version       display the version of links.ml and exit

TLS Options:
      --tls.enable    run tls server rather than standard http
  -c, --tls.cert=     path to ssl cert file
  -k, --tls.key=      path to ssl key file

Help Options:
  -h, --help          Show this help message
```

For example:

```
$ links -s "http://your-domain.com" -b "0.0.0.0:80" -d links.db
```

## Using as a library:

Links.ml also has a Go client library which you can use, which adds a simple
wrapper around an http call, to make shortening links simpler. Download it
using the following `go get` command:

```
$ go get -u github.com/lrstanley/links.ml/client
```

View the documentation [here](https://godoc.org/github.com/lrstanley/links.ml/client)

## Migrating

If you were using links.ml before it was rewritten in Go (and was using
MySQL), then this is how you can migrate the database:

```
$ links --migrate --migrate-info "user:passwd@tcp(your-server:3306)/links_db"
```

Depending on the amount of links within the old database, this may take a while
to complete. I would recommend backing things up prior just in case as well.

## Contributing

Below are a few guidelines if you would like to contribute. Keep the code
clean, standardized, and much of the quality should match Golang's standard
library and common idioms.

   * Always test using the latest Go version.
   * Always use `gofmt` before committing anything.
   * Always have proper documentation before committing.
   * Keep the same whitespacing, documentation, and newline format as the
     rest of the project.
   * Only use 3rd party libraries if necessary. If only a small portion of
     the library is needed, simply rewrite it within the library to prevent
     useless imports.
   * Also see [golang/go/wiki/CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments)

## API

Shortening a link is quite easy. simply send a `POST` request to `https://links.ml/add`,
which will return JSON-safe information as shown below:

```
$ curl --data "url=http://google.com" https://links.ml/add
{"url": "https://links.ml/27f4", "success": true}
```

#### Password protection

You can also password protect a link, simply by adding a `password` variable to the payload:

```
$ curl --data 'url=https://google.com/example&encrypt=y0urp4$$w0rd' https://links.ml/add
{"url": "https://links.ml/abc123", "success": true}
```

## License

```
LICENSE: The MIT License (MIT)
Copyright (c) Liam Stanley <me@liamstanley.io>

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
