// @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/golang/src/example-api-integration/hello-server-unprotected.js

package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type SuccessResponse struct {
    Message string `json:"message"`
}

func helloHandler(response http.ResponseWriter, request *http.Request) {

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusOK)

    json.NewEncoder(response).Encode(SuccessResponse{Message: "Hello, World!"})

    log.Println("Hello response sent...")
}

func main() {
    http.HandleFunc("/", helloHandler)

    log.Println("Server listening on http://localhost:8002")
    http.ListenAndServe(":8002", nil)
}
