// @link https://approov.io/docs/v2.0/approov-usage-documentation/#example-api-integration
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/golang/src/example-api-integration/hello-server-protected.js

package main

import (
    "fmt"
    "encoding/json"
    "encoding/base64"
    "log"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
)

const Log_Tag = "| APPROOV:"
const base64Secret = "h+CX0tOzdAAR9l15bWAqvq7w9olk66daIH+Xk+IAHhVVHszjDzeGobzNnqyRze3lw/WVyWrc2gZfh3XXfBOmww=="

type SuccessResponse struct {
    Message string `json:"message"`
}

type BadRequest struct {
    Error string `json:"error"`
}

func logApproov(message, log_level string) {
    log.Println(log_level, Log_Tag, message)
}

func sendBadRequestResponse(response http.ResponseWriter, message string) {

    logApproov(message, "ERROR")

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusBadRequest)

    json.NewEncoder(response).Encode(BadRequest{Error: "Bad Request"})
}

func hello(response http.ResponseWriter, request *http.Request) {

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusOK)

    json.NewEncoder(response).Encode(SuccessResponse{Message: "Hello, World!"})

    log.Println("Hello response sent...")
}

func verifyApproovToken(response http.ResponseWriter, request *http.Request)  (*jwt.Token, error) {

    approovToken := request.Header["Approov-Token"]

    if approovToken == nil {
        return nil, fmt.Errorf("token is missing")
    }

    token, err := jwt.Parse(approovToken[0], func(token *jwt.Token) (interface{}, error) {

        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Signing method mismatch.")
        }

        secret, err := base64.StdEncoding.DecodeString(base64Secret)

        if err != nil {
            return nil, fmt.Errorf(err.Error())
        }

        return secret, nil
    })

    return token, err
}

func checkApproovToken(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

    return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

        token, err := verifyApproovToken(response, request)

        if err != nil {
            sendBadRequestResponse(response, err.Error())
            return
        }

        if ! token.Valid {
            sendBadRequestResponse(response, "the token is invalid.")
            return
        }

        logApproov("valid token", "INFO")

        endpoint(response, request)
    })
}

func main() {
    http.Handle("/", checkApproovToken(hello))

    log.Println("Server listening on http://localhost:8002")
    http.ListenAndServe(":8002", nil)
}
