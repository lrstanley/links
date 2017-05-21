# Links.ml
### URL Shortener

## API usage

#### Basic

Shortening a link is quite easy. simply send a `POST` request to `https://links.ml/add`, which will return JSON-safe information as shown below:

```
$ curl --data "url=http://google.com" https://links.ml/add
{"url": "https://links.ml/27f4", "success": true}
```

#### Password protection

You can also password protect a link, simply by adding a `password` variable to the payload:

```
$ curl --data "url=http://google.com&password=test" https://links.ml/add
{"url": "https://links.ml/3BDd", "success": true}
```

If you happen to have any further issues, feel free to submit it with more
information [here](https://github.com/lrstanley/links.ml/issues)
