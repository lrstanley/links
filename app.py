#!/usr/bin/python
import flask
import os
import time
import json
import redis
import re

app = flask.Flask(__name__)


redis_prefix = 'links.ml'


@app.route('/')
@app.route('/<page>')
def main(page=None):
    if not page:
        page = 'index'
        tmp = '/'
    else:
        tmp = '/' + page
    if os.path.isfile('templates/%s.html' % page):
        return flask.render_template(page + '.html', page=tmp)
    return flask.abort(404)


@app.route('/add', methods=['POST'])
def add():
    form = list(flask.request.form)
    if len(form) != 1:
        data = {
            'success': False,
            'message':
            'Please enter a URL'
        }
        return json.dumps(data)
    url = form[0]
    if not is_url(url):
        data = {
            'success': False,
            'message': 'Please enter a valid URL'
        }
        return json.dumps(data)

    data = {
        'url': url,
        'success': True,
        'message': """
            Successfully shortened your URL!
            <input type="text" value="%s" class="form-control tip"
                   id="url-to-copy" title="Press CTRL+C to copy" readonly>"""
    }
    data['message'] = data['message'] % data['url']
    return json.dumps(data)


def is_url(url):
    regex = re.compile(
        r'^(?:http|ftp)s?://' # http:// or https://
        r'(?:(?:[A-Z0-9](?:[A-Z0-9-]{0,61}[A-Z0-9])?\.)+(?:[A-Z]{2,6}\.?|[A-Z0-9-]{2,}\.?)|' # domain...
        r'\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|' # ...or ipv4
        r'\[?[A-F0-9]*:[A-F0-9:]+\]?)' # ...or ipv6
        r'(?::\d+)?' # optional port
        r'(?:/?|[/?]\S+)$', re.IGNORECASE)
    if re.match(regex, url):
        return True
    else:
        return False


@app.context_processor
def utility_processor():
    def example():
        return 'yolo'

    return dict(
        example=example
    )


@app.errorhandler(404)
def page_not_found(error):
    return flask.render_template('404.html'), 404

 
def check_db():
    # Done on website-initialization
    pass


# Debug should normally be false, so we don't display hazardous information!
db = redis.StrictRedis(host='localhost', port=6379, db=0)
check_db()
if __name__ == '__main__':
    app.debug = True
    app.run(host='0.0.0.0', port=5000)



# ## Examples of stuff ##

# # Global
# # Loaded from DB on init
# # Save on change and reload global variable
# dadb = {
#     'UUID': {
#         'visit_count': 20,
#         'time': 1230912390,
#         'url': 'http://google.com'
#     },
#     'UUID2': {
#         'visit_count': 220,
#         'time': 1230922210,
#         'url': 'http://yolo.com'
#     }
# }

# # Global
# # Loaded from DB on init
# # Add URL every time one has been created
# # Save globally but not to DB
# # Adds:
# #   - Faster searching to see if a URL already exists
# db_urls = {
#     'http://yolo.com': 'UUID2',
#     'http://google.com': 'UUID'
# }


# # Global
# # Loaded from DB on init
# # global-page-counter and per-UUID counter
# # Change on page load of selected page
# # Save to global and db on each change
# count += 1
# dadb['UUID']['visit_count'] += 1