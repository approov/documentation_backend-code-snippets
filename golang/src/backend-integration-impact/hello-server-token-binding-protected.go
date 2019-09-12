// @link https://approov.io/docs/v2.0/approov-usage-documentation/#backend-integration-impact
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/golang/src/backend-integration-impact/hello-server-token-binding-protected.js

package main

import (
    "fmt"
    "encoding/json"
    "encoding/base64"
    "log"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
    "crypto/sha256"
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

func sendHelloRequestResponse(response http.ResponseWriter, request *http.Request) {

    response.Header().Set("Content-Type", "application/json")
    response.WriteHeader(http.StatusOK)

    json.NewEncoder(response).Encode(SuccessResponse{Message: "Hello, World!"})

    log.Println("Hello response sent...")
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

func verifyApproovTokenBinding(token *jwt.Token, request *http.Request) (jwt.Claims, error) {

    claims := token.Claims

    token_binding_payload, has_pay_key := claims.(jwt.MapClaims)["pay"]

    if ! has_pay_key {
        // We log a warning and don't return an error, because the Approov fail-over server doesn't provide the `pay` key.
        logApproov("Key 'pay' in the decoded token is missing or empty.", "WARNING")
        return claims, nil
    }

    // We use the Authorization token here, but feel free to use another header. However, you need to bind this header to the
    // Approov token in the mobile app.
    token_binding_header := request.Header["Authorization"]

    if len(token_binding_header) != 1 {
        return claims, fmt.Errorf("The header to perform the verification for the token binding is missing, empty or has more than one entry.")
    }

    // We need to hash and base64 encode the token binding header, because we need to compare it in the same way it was
    // included in the Approov token on the mobile app.
    token_binding_header_hashed := sha256.Sum256([]byte(token_binding_header[0]))
    token_binding_header_encoded := base64.StdEncoding.EncodeToString(token_binding_header_hashed[:])

    if token_binding_payload != token_binding_header_encoded {
        return claims, fmt.Errorf("Invalid token binding.")
    }

    return claims, nil
}

func checkApproovTokenBinding(next http.Handler) http.Handler {

    return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

        token, err := verifyApproovToken(response, request)

        if err != nil {
            sendBadRequestResponse(response, err.Error())
            return
        }

        if ! token.Valid {
            sendBadRequestResponse(response, "the token is invalid")
            return
        }

        logApproov("Valid token.", "INFO")

        claims, err := verifyApproovTokenBinding(token, request)

        if err != nil {
            sendBadRequestResponse(response, err.Error())
            return
        }

        // Use the returned claims to perform any additional checks.
        _ = claims

        logApproov("Valid token binding.", "INFO")

        next.ServeHTTP(response, request)
    })
}

func main() {
    helloHandler := http.HandlerFunc(sendHelloRequestResponse)
    http.Handle("/", checkApproovTokenBinding(helloHandler))

    log.Println("Server listening on http://localhost:8002")
    http.ListenAndServe(":8002", nil)
}
