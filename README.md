links.ml
========

Links.ml - URL Shortener

Installation Instructions:
--------------------------

 - `$ sudo apt-get install python-pip redis-server git-core`
 - `$ sudo pip install flask redis`
 - `$ git clone https://github.com/Liamraystanley/links.ml.git`
 - `$ cd links.ml`
 - Edit the very top of `app.py` to customize your settings (E.g. `$ nano app.py`)
 - `$ python app.py`
 - Head to `http://localhost:5000/`

    Note: If you want to access links.ml without a port, you must use Nginx or Apache http-proxies to forward over to port 80
