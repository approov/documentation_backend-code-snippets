// @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/nodejs/src/example-api-integration/hello-server-unprotected.js

const express = require('express')
const app = express()

// simple 'hello world' endpoint.
app.get('/', function (req, res, next) {
    res.json({message: "Hello World!"})
})

// Create and run the HTTP server
app.listen(8002, function () {
  console.debug("Server listening on %s", "http://localhost:8002")
})
