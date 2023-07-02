<!-- template:define:options
{
  "nodescription": true
}
-->
![logo](https://liam.sh/-/gh/svg/lrstanley/links?icon=ph%3Alink-simple&icon.height=65&layout=left&icon.color=rgba%280%2C+184%2C+126%2C+1%29&font=1.3)

<!-- template:begin:header -->
<!-- do not edit anything in this "template" block, its auto-generated -->

<p align="center">
  <a href="https://github.com/lrstanley/links/releases">
    <img title="Release Downloads" src="https://img.shields.io/github/downloads/lrstanley/links/total?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/tags">
    <img title="Latest Semver Tag" src="https://img.shields.io/github/v/tag/lrstanley/links?style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/commits/master">
    <img title="Last commit" src="https://img.shields.io/github/last-commit/lrstanley/links?style=flat-square">
  </a>




  <a href="https://github.com/lrstanley/links/actions?query=workflow%3Atest+event%3Apush">
    <img title="GitHub Workflow Status (test @ master)" src="https://img.shields.io/github/actions/workflow/status/lrstanley/links/test.yml?branch=master&label=test&style=flat-square">
  </a>

  <a href="https://codecov.io/gh/lrstanley/links">
    <img title="Code Coverage" src="https://img.shields.io/codecov/c/github/lrstanley/links/master?style=flat-square">
  </a>

  <a href="https://pkg.go.dev/github.com/lrstanley/links">
    <img title="Go Documentation" src="https://pkg.go.dev/badge/github.com/lrstanley/links?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/lrstanley/links">
    <img title="Go Report Card" src="https://goreportcard.com/badge/github.com/lrstanley/links?style=flat-square">
  </a>
</p>
<p align="center">
  <a href="https://github.com/lrstanley/links/issues?q=is:open+is:issue+label:bug">
    <img title="Bug reports" src="https://img.shields.io/github/issues/lrstanley/links/bug?label=issues&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/issues?q=is:open+is:issue+label:enhancement">
    <img title="Feature requests" src="https://img.shields.io/github/issues/lrstanley/links/enhancement?label=feature%20requests&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/pulls">
    <img title="Open Pull Requests" src="https://img.shields.io/github/issues-pr/lrstanley/links?label=prs&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/releases">
    <img title="Latest Semver Release" src="https://img.shields.io/github/v/release/lrstanley/links?style=flat-square">
    <img title="Latest Release Date" src="https://img.shields.io/github/release-date/lrstanley/links?label=date&style=flat-square">
  </a>
  <a href="https://github.com/lrstanley/links/discussions/new?category=q-a">
    <img title="Ask a Question" src="https://img.shields.io/badge/support-ask_a_question!-blue?style=flat-square">
  </a>
  <a href="https://liam.sh/chat"><img src="https://img.shields.io/badge/discord-bytecord-blue.svg?style=flat-square" title="Discord Chat"></a>
</p>
<!-- template:end:header -->

<!-- template:begin:toc -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :link: Table of Contents

  - [💻 Installation](#computer-installation)
    - [Container Images (ghcr)](#whale-container-images-ghcr)
    - [🧰 Source](#toolbox-source)
  - [⚙️ Usage](#gear-usage)
    - [Dot-env](#dot-env)
    - [Example](#example)
    - [Google SafeBrowsing](#google-safebrowsing)
  - [Using as a library](#using-as-a-library)
    - [Example](#example-1)
  - [API](#api)
      - [Password protection](#password-protection)
  - [Support &amp; Assistance](#raising_hand_man-support--assistance)
  - [🤝 Contributing](#handshake-contributing)
  - [License](#balance_scale-license)
<!-- template:end:toc -->

## :computer: Installation

Check out the [releases](https://github.com/lrstanley/links/releases)
page for prebuilt versions.

<!-- template:begin:ghcr -->
<!-- do not edit anything in this "template" block, its auto-generated -->
### :whale: Container Images (ghcr)

```console
$ docker run -it --rm ghcr.io/lrstanley/links:master
$ docker run -it --rm ghcr.io/lrstanley/links:0.8.0
$ docker run -it --rm ghcr.io/lrstanley/links:latest
$ docker run -it --rm ghcr.io/lrstanley/links:0.7.0
$ docker run -it --rm ghcr.io/lrstanley/links:0.6.0
```
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

	fmt.Printf("shortened: %s
", uri.String())
}
```

## API

Shortening a link is quite easy. simply send a `POST` request to `https://example.com/add`,
which will return JSON-safe information as shown below:

```
$ curl --data "url=http://google.com" https://example.com/add
{"url": "https://example.com/27f4", "success": true}
```

#### Password protection

You can also password protect a link, simply by adding a `password` variable to the payload:

```
$ curl --data 'url=https://google.com/example&encrypt=y0urp4$$w0rd' https://example.com/add
{"url": "https://example.com/abc123", "success": true}
```

<!-- template:begin:support -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :raising_hand_man: Support & Assistance

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for
     guidelines on ensuring everyone has the best experience interacting with
     the community.
* :raising_hand_man: Take a look at the [support](.github/SUPPORT.md) document on
     guidelines for tips on how to ask the right questions.
* :lady_beetle: For all features/bugs/issues/questions/etc, [head over here](https://github.com/lrstanley/links/issues/new/choose).
<!-- template:end:support -->

<!-- template:begin:contributing -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :handshake: Contributing

* :heart: Please review the [Code of Conduct](.github/CODE_OF_CONDUCT.md) for guidelines
     on ensuring everyone has the best experience interacting with the
    community.
* :clipboard: Please review the [contributing](.github/CONTRIBUTING.md) doc for submitting
     issues/a guide on submitting pull requests and helping out.
* :old_key: For anything security related, please review this repositories [security policy](https://github.com/lrstanley/links/security/policy).
<!-- template:end:contributing -->

<!-- template:begin:license -->
<!-- do not edit anything in this "template" block, its auto-generated -->
## :balance_scale: License

```
MIT License

Copyright (c) 2014 Liam Stanley <me@liamstanley.io>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

_Also located [here](LICENSE)_
<!-- template:end:license -->
