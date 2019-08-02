# @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
# @link https://github.com/approov/documentation_backend-code-snippets/blob/master/python/src/example-api-integration/hello-server-unprotected.py

from flask import Flask, jsonify

app = Flask(__name__)

@app.route("/")
def hello():
   return jsonify({"message": "Hello World!"})

if __name__ == "__main__":
   app.run()
