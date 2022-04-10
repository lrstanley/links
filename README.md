<!-- template:begin:header -->
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- template:end:toc -->

## :computer: Installation

Check out the [releases](https://github.com/lrstanley/links/releases)
page for prebuilt versions.

<!-- template:begin:ghcr -->
<!-- template:end:ghcr -->

### :toolbox: Source

Note that you must have [Go](https://golang.org/doc/install) installed (latest is usually best).

```console
$ git clone https://github.com/lrstanley/links.git && cd links
$ make
$ ./links --help
```

## :gear: Usage

```
$ links --help
Usage:
  links [OPTIONS] [add | delete]

Application Options:
  -s, --site-name=                      site url, used for url generation (default: https://links.wtf) [$SITE_URL]
      --session-dir=                    optional location to store temporary sessions [$SESSION_DIR]
  -q, --quiet                           don't log to stdout [$QUIET]
      --debug                           enable debugging (pprof endpoints) [$DEBUG]
  -b, --http=                           ip:port pair to bind to (default: :8080) [$HTTP]
  -p, --behind-proxy                    if X-Forwarded-For headers should be trusted [$PROXY]
  -d, --db=                             path to database file (default: store.db) [$DB_PATH]
      --key-length=                     default length of key (uuid) for generated urls (default: 4) [$KEY_LENGTH]
      --http-pre-include=               HTTP include which is included directly after css is included (near top of the page)
                                        [$HTTP_PRE_INCLUDE]
      --http-post-include=              HTTP include which is included directly after js is included (near bottom of the page)
                                        [$HTTP_POST_INCLUDE]
  -e, --export-file=                    file to export db to (default: links.export)
      --export-json                     export db to json elements
  -v, --version                         display the version of links.wtf and exit

TLS Options:
      --tls.enable                      run tls server rather than standard http [$TLS_ENABLE]
  -c, --tls.cert=                       path to ssl cert file [$TLS_CERT_PATH]
  -k, --tls.key=                        path to ssl key file [$TLS_KEY_PATH]

Safe Browsing Support:
      --safebrowsing.api-key=           Google API Key used for querying SafeBrowsing, disabled if not provided (see:
                                        https://github.com/lrstanley/links#google-safebrowsing) [$SAFEBROWSING_API_KEY]
      --safebrowsing.db=                path to SafeBrowsing database file (default: safebrowsing.db) [$SAFEBROWSING_DB_PATH]
      --safebrowsing.update-period=     duration between updates to the SafeBrowsing API local database (default: 1h)
                                        [$SAFEBROWSING_UPDATE_PERIOD]
      --safebrowsing.redirect-fallback  if the SafeBrowsing request fails (local cache, and remote hit), this still lets the
                                        redirect happen [$SAFEBROWSING_REDIRECT_FALLBACK]

Prometheus Metrics:
      --prom.enabled                    enable exposing of prometheus metrics (on std port, or --prometheus.addr) [$PROM_ENABLED]
      --prom.addr=                      expose on custom address/port, e.g. ':9001' (all ips) or 'localhost:9001' (local only)
                                        [$PROM_ADDR]
      --prom.endpoint=                  endpoint to expose metrics on (default: /metrics) [$PROM_ENDPOINT]

Help Options:
  -h, --help                            Show this help message

Available commands:
  add     add a link
  delete  delete a link, id, or link matching an author

```

### Dot-env

Links also supports a `.env` file for loading environment variables. Example:

```
SAFEBROWSING_API_KEY=<your-secret-key>
FOO=BAR
```

### Example

```
$ links -s "http://your-domain.com" -b "0.0.0.0:80" -d links.db
```

### Google SafeBrowsing

Links supports utilizing [Google SafeBrowsing](https://safebrowsing.google.com/),
(see `Safe Browsing Support` under ##Usage). This helps prevent users being
redirected to malicious or otherwise harmful websites. It does require a Google
Developer account (free).

   1. Go to the [Google Developer Console](https://console.developers.google.com/)
   2. Create a new project (dropdown top left, click `NEW PROJECT`)
   3. Enable the "Safe Browsing" API
   4. Create a new credential (API key)

Screenshot example:
![](https://i.imgur.com/VHyPYqi.png)

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

<!-- template:begin:support -->
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- template:end:license -->
