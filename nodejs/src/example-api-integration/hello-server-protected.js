// @link https://approov.io/docs/v2.1/approov-usage-documentation/#example-api-integration
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/nodejs/src/example-api-integration/hello-server-protected.js

const express = require('express')
const jwt = require('express-jwt')
const app = express()

const dotenv = require('dotenv').config()

if (dotenv.error) {
  console.debug('FAILED TO PARSE `.env` FILE | ' + dotenv.error)
}

const SECRET = dotenv.parsed.APPROOV_BASE64_SECRET

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
    res.status(401)
    res.json(ERROR_RESPONSE_BODY)
    return
  }

  next()
  return
}

// Callback to handles when an Approov token is successfully validated.
const handlesApproovTokenSuccess = function(req, res, next) {
  if (req.approovTokenError === true) {
    res.status(401)
    res.json(ERROR_RESPONSE_BODY)
    return
  }

  next()
  return
}

// Intercepts all calls to the 'hello world' endpoint and validates the Approov token.
app.use('/', checkApproovToken)

// Handles failure in validating the Approov token
app.use('/', handlesApproovTokenError)

// Handles requests where the Approov token is a valid one.
app.use('/', handlesApproovTokenSuccess)

// simple 'hello world' endpoint.
app.get('/', function (req, res, next) {
    res.json({message: "Hello Approov token protected server!"})
})

// Create and run the HTTP server
app.listen(8002, function () {
  console.debug("Server listening on %s", "http://localhost:8002")
})
