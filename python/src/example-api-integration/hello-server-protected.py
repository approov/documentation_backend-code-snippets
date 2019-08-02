# @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
# @link https://github.com/approov/documentation_backend-code-snippets/blob/master/python/src/example-api-integration/hello-server-protected.py

from flask import Flask, request, abort, jsonify
import jwt # https://github.com/jpadilla/pyjwt/
import json
import base64

app = Flask(__name__)

# Token secret value obtained with the Approov CLI tool: approov secret <admin.tok> -get
SECRET = bytes("VcXl2IgVUgzYZTHloVJ/gmkHyUuvXIRrN+aCjk7/4ZPwq24/XSF1DaTOcHCX6j3eUNUX3i4GBfcpNZKSeUefqA==","ascii")

# Function to check the validity of the token
def verifyToken(token):
  try:
    # Decode our token, allowing only the HS256 algorithm, using our base64 encoded SECRET
    tokenContents = jwt.decode(token, base64.b64decode(SECRET), algorithms=['HS256'])
    return tokenContents
  except jwt.ExpiredSignatureError as e:
    # Signature has expired, token is bad
    return None
  except jwt.InvalidTokenError as e:
    # Token could not be decoded, token is bad
    return None

@app.route("/")
def hello():
  # Get the Approov Token from header
  token = request.headers.get("Approov-Token")

  # If we didn't find a token, then reject the request
  if token == "":
    abort(401)

  tokenContents = verifyToken(token)

  if (tokenContents == None):
    abort(401)

  return jsonify({"message": "Hello World!"})

if __name__ == "__main__":
   app.run()
