#!/usr/bin/python
import flask
from time import time as timenow
from utils import *
from thread import start_new_thread as bg

app = flask.Flask(__name__)


@app.route("/")
def main():
    # We don't need to increment counters in the foreground
    bg(counter, (None,))
    return flask.render_template("index.html")


@app.route("/<link>")
def route(link=None):
    if not link:
        return flask.redirect("/")
    if not valid_uri(link):
        return flask.redirect("/")

    # We don't need to increment counters in the foreground
    bg(counter, (None,))

    link = table("urls").find_one(id=link)

    # If it doesn't exist
    if not link:
        return flask.redirect("/")

    # Increment it in the background
    bg(counter, (link.id, link.hits))

    # If it's password protected
    if link.password:
        return flask.render_template("password.html")

    return flask.redirect(link.url)


@app.route("/decrypt", methods=["POST"])
def decrypt():
    required = [
        {"name": "path", "error": "No ID supplied"},
        {"name": "password", "error": "A password is required"}
    ]
    form = flask.request.get_json() if flask.request.get_json() else flask.request.form

    for item in required:
        if item["name"] not in form:
            return err(item["error"])
    link = form["path"].strip("/")
    if not valid_uri(link):
        return err("Insufficient or malformed ID")
    link = table("urls").find_one(id=link)

    # If it doesn't actually exist
    if not link:
        return err("ID provided does not exist")

    if hash(form["password"]) != link.password:
        return err("Incorrect password entered")
    return success(url=link.url)


@app.route("/add", methods=["POST"])
def add():
    max_uuid_fails = 20
    form = flask.request.get_json() if flask.request.get_json() else flask.request.form
    if "url" not in form:
        return err("Please enter a valid URL")
    url = form['url'].strip()
    passwd = form["password"] if "password" in form else None
    passwd = hash(passwd) if passwd and len(passwd) > 0 else None
    if not valid_url(url):
        return err("Please enter a valid URL")

    # First, see if it's already in the database to prevent dups
    # Also note: Passworded links are unique
    link = table("urls").find_one(url=url, password=passwd)

    if link:
        return success(url=cnf("General", "url", "https://links.ml").rstrip("/") + "/" + link.id)

    _uuid_count = 0
    try:
        current_uuids = [x.id for x in list(db().query("SELECT id FROM urls"))]
    except:
        current_uuids = []
    while True:
        _uuid_count += 1
        if _uuid_count > max_uuid_fails:
            # It's been checked too many times, and performance is likely degrading
            msg = "ERROR: While attempting to generate a new UUID, it failed %s times. If this message re-occurs, this is likely a sign that key-generation is failing. Consult the main.cfg file."
            print(msg % str(max_uuid_fails))
            return err("Failed to generate UUID")
        uuid = new_uuid()
        if uuid not in current_uuids:
            table("urls").insert(dict(id=uuid, url=url, hits=0, password=passwd, time=int(timenow())))
            return success(url=cnf("General", "url", "https://links.ml").rstrip("/") + "/" + uuid)


@app.route("/stats")
def stats():
    count = table("stats").find_one(id="hits")
    count = int(count.count) if count else 1
    shortened = sum([int(link.hits) for link in list(table("urls").all())])
    return flask.jsonify({
        "success": True,
        "views": {"raw": count, "clean": "{:,d}".format(count)},
        "links": {"raw": shortened, "clean": "{:,d}".format(shortened)},
    })


@app.errorhandler(404)
def page_not_found(error):
    return flask.render_template("404.html"), 404


if __name__ == "__main__":
    app.debug = True
    host = cnf("General", "bindhost", "0.0.0.0")
    port = int(cnf("General", "port", 3333))
    app.run(host=host, port=port)
