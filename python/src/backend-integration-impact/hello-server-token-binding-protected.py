# @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
# @link https://github.com/approov/documentation_backend-code-snippets/blob/master/python/src/backend-integration-impact/hello-server-token-binding-protected.py

from flask import Flask, request, abort, jsonify, make_response
import jwt # https://github.com/jpadilla/pyjwt/
import json
import base64
import hashlib

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
    app.logger.info("Missing the Approov token header.")
    return None

  try:
    # Decode our token, allowing only the HS256 algorithm, using our base64 encoded SECRET
    approov_token_claims = jwt.decode(approov_token, base64.b64decode(SECRET), algorithms=['HS256'])
  except jwt.ExpiredSignatureError as e:
    # Signature has expired, token is bad
    app.logger.info(e)
    return None
  except jwt.InvalidTokenError as e:
    # Token could not be decoded, token is bad
    app.logger.info(e)
    return None

  return approov_token_claims

def verifyApproovTokenBinding(request, approov_token_claims):

  if approov_token_claims is None:
    app.logger.info("Approov token claims cannot be empty.")
    return False

  if not 'pay' in approov_token_claims:
    app.logger.info("Key 'pay' in the decoded Approov token is missing. ")

    # We need to return True because if the Approov token is coming from the fail-over server it will not contain it.
    return None

  token_binding_header = request.headers.get("Authorization")

  if not token_binding_header:
    app.logger.info("Missing the token binding header in this request.")
    return False

  # We need to hash and base64 encode the token binding header, because that's how it was included in the Approov
  # token on the mobile app.
  token_binding_header_hash = hashlib.sha256(token_binding_header.encode('utf-8')).digest()
  token_binding_header_encoded = base64.b64encode(token_binding_header_hash).decode('utf-8')

  if approov_token_claims['pay'] == token_binding_header_encoded:
    app.logger.info("Valid token binding.")
    return True

  app.logger.info("Approov token binding claim not matching token binding in header.")
  return False

@app.route("/")
def hello():

  approov_token_claims = verifyApproovToken(request)

  if approov_token_claims is None:
    app.logger.info("Invalid Approov token.")
    abort(make_response(jsonify({}), 401))

  if verifyApproovTokenBinding(request, approov_token_claims) == False:
    app.logger.info("Invalid token binding.")
    abort(make_response(jsonify({}), 401))

  return jsonify({"message": "Hello World!"})


if __name__ == "__main__":
   app.run()
