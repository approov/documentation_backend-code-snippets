# @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
# @link https://github.com/approov/documentation_backend-code-snippets/blob/master/python/src/example-api-integration/hello-server-protected.py

from flask import Flask, request, abort, jsonify, make_response
import jwt # https://github.com/jpadilla/pyjwt/
import json
import base64

from os import getenv
from dotenv import load_dotenv, find_dotenv
load_dotenv(find_dotenv(), override=True)

# Token secret value obtained with the Approov CLI tool: approov secret <admin.tok> -get
SECRET = getenv('APPROOV_BASE64_SECRET')

app = Flask(__name__)

def verifyApproovToken(request):
  # Get the Approov Token from header
  approov_token = request.headers.get("Approov-Token")

  # If we didn't find a token, then reject the request
  if approov_token == "":
    return None

  try:
    # Decode our token, allowing only the HS256 algorithm, using our base64 encoded SECRET
    approov_token_claims = jwt.decode(approov_token, base64.b64decode(SECRET), algorithms=['HS256'])
  except jwt.ExpiredSignatureError as e:
    # Signature has expired, token is bad
    return None
  except jwt.InvalidTokenError as e:
    # Token could not be decoded, token is bad
    return None

  return approov_token_claims

@app.route("/")
def hello():
  approov_token_claims = verifyApproovToken(request)

  if approov_token_claims is None:
    abort(make_response(jsonify({}), 401))

  return jsonify({"message": "Hello World!"})

if __name__ == "__main__":
   app.run()
