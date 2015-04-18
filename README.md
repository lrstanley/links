links.ml
========

Links.ml - URL Shortener

Installation Instructions:
--------------------------

 - `$ sudo apt-get install python-pip redis-server git-core`
 - `$ sudo pip install flask redis`
 - `$ git clone https://github.com/Liamraystanley/links.ml.git`
 - `$ cd links.ml`
 - Edit the very top of `app.py` to customize your settings (E.g. `$ vim app.py`)
 - `$ python app.py`
 - Head to `http://localhost:5000/`

    Note: If you want to access links.ml without a port, you must use Nginx or Apache http-proxies to forward over to port 80


API usage
---------

#### Basic

Shortening a link is quite easy. simply send a `POST` request to `https://links.ml/add`, which will return JSON-safe information as shown below:

```
$ curl --data "url=http://google.com" https://links.ml/add
{"url": "https://links.ml/BA2b", "success": true}
```

#### Password protection

You can also password protect a link, simply by adding a `password` variable to the payload:

```
$ curl --data "url=http://google.com&password=test" https://links.ml/add
{"url": "https://links.ml/98c1", "success": true}
```
