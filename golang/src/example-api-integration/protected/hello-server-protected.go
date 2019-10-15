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

const base64Secret = "h+CX0tOzdAAR9l15bWAqvq7w9olk66daIH+Xk+IAHhVVHszjDzeGobzNnqyRze3lw/WVyWrc2gZfh3XXfBOmww=="

type SuccessResponse struct {
    Message string `json:"message"`
}

type ErrorResponse struct {}

func errorResponse(response http.ResponseWriter, statusCode int, message string) {
    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(statusCode)
    json.NewEncoder(response).Encode(ErrorResponse{})
}

func helloHandler(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(SuccessResponse{Message: "Hello, World!"})
}

func verifyApproovToken(response http.ResponseWriter, request *http.Request)  (*jwt.Token, error) {
    approovToken := request.Header["Approov-Token"]

    if approovToken == nil {
        return nil, fmt.Errorf("Token is missing.")
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

func makeApproovCheckerHandler(handler func(http.ResponseWriter, *http.Request)) http.Handler {
    return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

        token, err := verifyApproovToken(response, request)

        if err != nil {
            errorResponse(response, http.StatusUnauthorized, err.Error())
            return
        }

        if ! token.Valid {
            errorResponse(response, http.StatusUnauthorized, "Token is invalid.")
            return
        }

        handler(response, request)
    })
}

func main() {
    http.Handle("/", makeApproovCheckerHandler(helloHandler))
    log.Println("Server listening on http://localhost:8002")
    http.ListenAndServe(":8002", nil)
}
