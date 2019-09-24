// @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/nodejs/src/backend-integration-impact/hello-server-token-binding-protected.js

const express = require('express')
const jwt = require('express-jwt')
const crypto = require('crypto')
const app = express()

const dotenv = require('dotenv').config()

if (dotenv.error) {
  console.debug('FAILED TO PARSE `.env` FILE | ' + dotenv.error)
}

const SECRET = dotenv.parsed.APPROOV_BASE64_SECRET

const isEmpty = function(value) {
  return  (value === undefined) || (value === null) || (value === '')
}

const isString = function(value) {
  return (typeof(value) === 'string')
}

const isEmptyString = function(value) {
  return (isEmpty(value) === true) || (isString(value) === false) ||  (value.trim() === '')
}

const ERROR_RESPONSE_BODY = {}

// Callback that performs the Approov token check using the express-jwt library
const checkApproovToken = jwt({
  secret: Buffer.from(SECRET, 'base64'), // decodes the Approov secret
  requestProperty: 'approovTokenDecoded',
  getToken: function fromApproovTokenHeader(req) {
    req.approovTokenError = false
    return req.get('Approov-Token')
  },
  algorithms: ['HS256']
})

// Callback to handle the errors occurred while checking the Approov token.
const handlesApproovTokenError = function(err, req, res, next) {

  if (err.name === 'UnauthorizedError') {
    req.approovTokenError = true

    console.debug('APPROOV TOKEN ERROR: %s', err)

    res.status(401)
    res.json(ERROR_RESPONSE_BODY)
    return
  }

  next()
  return
}

// Callback to handles when an Approov token is successfully validated.
const handlesApproovTokenSuccess = function(req, res, next) {
  if (req.approovTokenError === false) {
    console.debug('VALID APPROOV TOKEN')
  }

  next()
  return
}

/**
 * @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
 * @link https://github.com/approov/documentation_backend-code-snippets/blob/master/nodejs/src/backend-integration-impact/hello-server-token-binding-protected.js
 */
const handlesApproovTokenBindingVerification = function(req, res, next) {

    // The decoded Approov token was added to the request object when the checked it at `checkApproovToken()`
    token_binding_payload = req.approovTokenDecoded.pay

    if (token_binding_payload === undefined) {
        console.debug("APPROOV TOKEN BINDING WARNING: key 'pay' is missing.")

        // We let the request to continue as usual, because the Approov fail-over server doesn't provide the `pay` key.
        next()
        return
    }

    if (isEmptyString(token_binding_payload)) {
        console.debug("APPROOV TOKEN BINDING ERROR: key 'pay' in the decoded token is missing or empty.")
        res.status(401)
        res.json(ERROR_RESPONSE_BODY)
        return
    }

    // We use here the Authorization token, but feel free to use another header, but you need to bind this  header to
    // the Approov token in the mobile app.
    token_binding_header = req.get('Authorization')

    if (isEmptyString(token_binding_header)) {
        console.debug("APPROOV TOKEN BINDING ERROR: Missing or empty header to perform the verification for the token binding.")
        res.status(401)
        res.json(ERROR_RESPONSE_BODY)
        return
    }

    // We need to hash and base64 encode the token binding header, because thats how it was included in the Approov
    // token on the mobile app.
    const token_binding_header_encoded = crypto.createHash('sha256').update(token_binding_header, 'utf-8').digest('base64')

    if (token_binding_payload !== token_binding_header_encoded) {
        console.debug("APPROOV TOKEN BINDING ERROR: Invalid token binding.")
        res.status(401)
        res.json(ERROR_RESPONSE_BODY)
        return
    }

    console.debug("VALID APPROOV TOKEN BINDING.")

    // Let the request continue as usual.
    next()
    return
}

// Intercepts all calls to the 'hello world' endpoint and validates the Approov token.
app.use('/', checkApproovToken)

// Handles failure in validating the Approov token
app.use('/', handlesApproovTokenError)

// Handles requests where the Approov token is a valid one.
app.use('/', handlesApproovTokenSuccess)

// Handles requests where the Approov token contains a bind to an another header in the request
app.use('/', handlesApproovTokenBindingVerification)

// simple 'hello world' endpoint.
app.get('/', function (req, res, next) {
    res.json({message: "Hello Approov token binding protected server!"})
})

// Create and run the HTTP server
app.listen(8002, function () {
  console.debug("Server listening on %s", "http://localhost:8002")
})
