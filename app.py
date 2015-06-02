#!/usr/bin/python
import flask
import os
import time
import json
import redis
import re
import string
import random

app = flask.Flask(__name__)

# Database prefix to use in redis to organize things
redis_prefix = 'links.ml'

# Prefix that the key's are appended to (needs the / at the end)
html_prefix = 'https://links.ml/'

# Length of the key you want. It's not incremental so
# it's recommended to have it be 4+ at least
key_length = 4

# If this is enabled it doesn't write to the DB, and resets on reboot
temporary_mode = False


@app.route('/')
@app.route('/<page>')
def main(page='index'):
    global db, urls, count
    # Up the global page count
    if os.path.isfile('templates/%s.html' % page):
        count['page'] += 1
        save('count', count)
        return flask.render_template(page + '.html')
    elif page in db:
        count['page'] += 1
        save('count', count)
        # Up the url count
        db[page]['visit_count'] += 1
        # Up the global url count
        count['urls'] += 1
        save('main', db)
        save('count', count)
        if 'password' in db[page]:
            return flask.render_template('password.html')
        return flask.redirect(db[page]['url'])
    else:
        return flask.redirect('/')


@app.route('/decrypt', methods=['POST'])
def decrypt():
    global db
    form = flask.request.form
    if flask.request.get_json():
        form = flask.request.get_json()

    if 'path' not in form:
        return json.dumps({'success': False, 'message': 'Insufficient or malformed path!'}), 400
    elif 'password' not in form:
        return json.dumps({'success': False, 'message': 'A password is required!'}), 400
    elif form['path'].strip('/') not in db:
        return json.dumps({'success': False, 'message': 'That URL doesn\'t exist!'}), 400
    elif form['password'] != db[form['path'].strip('/')]['password']:
        return json.dumps({'success': False, 'message': 'Incorrect password entered.'}), 400
    else:
        return json.dumps({'success': True, 'url': db[form['path'].strip('/')]['url']})


@app.route('/add', methods=['POST'])
def add():
    global db, urls
    form = flask.request.form
    # If headers are application/json, use this instead
    if flask.request.get_json():
        form = flask.request.get_json()
    if 'url' not in form:
        return json.dumps({'success': False, 'message': 'Please enter a valid URL'}), 400

    url = form['url'].strip()
    if not valid_url(url):
        return json.dumps({'success': False, 'message': 'Please enter a valid URL'}), 400

    passworded = False
    if 'password' in form:
        if len(form['password']) > 0:
            passworded = True

    # First, see if it's already in the DB to prevent duplications
    # Note: Passworded URL's are unique
    exists = False
    if url in urls and not passworded:
        exists = True

    if exists:
        data = {'url': html_prefix + urls[url], 'success': True}
    else:
        while True:
            tmp_uuid = new_uuid()
            if tmp_uuid not in db:
                break
        db[tmp_uuid] = {
            'visit_count': 0,
            'time': int(time.time()),
            'url': url
        }
        if passworded:
            db[tmp_uuid]['password'] = form['password']
        data = {'url': html_prefix + tmp_uuid, 'success': True}

        # Add it to the "urls" tmp variable, for faster duplication searching
        if not passworded:
            urls[url] = tmp_uuid
        save('main', db)
    return json.dumps(data), 200


def valid_url(url):
    if '//links.ml' in url.lower() or '//www.links.ml' in url.lower() or \
       len(url) < 1:
        return False
    if url.lower().startswith('www') or url[0].isdigit():
        url = 'http://' + url
    regex = re.compile(
        r'^http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+', re.IGNORECASE)
    if re.match(regex, url.strip()):
        return True
    else:
        return False

def new_uuid():
    return "".join([random.choice(string.hexdigits) for n in xrange(key_length)])


@app.route('/stats')
def stats():
    return flask.jsonify({
        'views': {
            'raw': count['page'],
            'clean': "{:,d}".format(count['page'])
        },
        'links': {
            'raw': count['urls'],
            'clean': "{:,d}".format(count['urls'])
        },
    })


@app.errorhandler(404)
def page_not_found(error):
    return flask.render_template('404.html'), 404


db, urls, count = {}, {}, {}
def check_db():
    # Done on website-initialization
    global db, urls, count
    if temporary_mode:
        db, urls, count = {}, {}, {'page': 0, 'urls': 0}
        return
    db = database.get(redis_prefix  + 'main')
    if not db:
        db = {}
    else:
        db = json.loads(db)
        for uuid in db:
            uuid_url = db[uuid]['url']
            urls[uuid_url] = uuid

    count = database.get(redis_prefix + 'count')
    if not count:
        count = {'page': 0, 'urls': 0}
    else:
        count = json.loads(count)


def save(name, data):
    if temporary_mode:
        return
    if isinstance(data, dict):
        data = json.dumps(data)

    database.set(redis_prefix + name, data)


database = redis.StrictRedis(host='localhost', port=6379, db=0)
check_db()
if __name__ == '__main__':
    app.debug = True
    app.run(host='0.0.0.0', port=5003)
