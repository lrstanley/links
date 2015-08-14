# Links.ml
### URL Shortener

## Setup

#### Databases

Create the database to use (make sure this matches the one listed within your
configuration file) by opening a mysql command prompt `$ mysql -u root -p`:

```sql
mysql> CREATE DATABASE links_db;
```

Generate a user to use for the database:

```sql
CREATE USER 'links_user'@'localhost' IDENTIFIED BY 'yourpassword';
```

And last but not least, grant that user permissions for the database:

```sql
mysql> GRANT ALL PRIVILEGES ON links_db.* TO 'links_user'@'localhost';
```

```sql
mysql> FLUSH PRIVILEGES;
```

Now, you can copy `example.cfg` to `main.cfg` and modify it to match the
above (alternatively, you can use `vim` instead of `nano`):

```
$ cp example.cfg main.cfg
$ nano main.cfg
```

#### The application

Install **pip** and **git**:

```
$ sudo apt-get install python-pip git-core
```

Install necessary Python libs:

```
$ sudo pip install mysql-python gunicorn dataset stuf
```

If you happen to run into an error when installing **mysql-python**, closely
related to `EnvironmentError: mysql_config not found`, and you're running Debian, execute the below command. Once done, re-run the above command.

```
$ sudo apt-get install libmysqlclient-dev python-dev
```

Once everythings installed, you should be able to download the latest **links.ml**:

```
$ git clone https://github.com/Liamraystanley/links.ml.git
$ cd links.ml
```

And now it should be able to run:

```
$ gunicorn -w 1 -b 0.0.0.0:yourport app:app
```

**Note:** gunicorn is just an example **WSGI** server, simply replace it with your favorite.


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
information [here](https://github.com/Liamraystanley/links.ml/issues)
