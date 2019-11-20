// @link https://approov.io/docs/v2.1/approov-usage-documentation/#backend-integration-impact
// @link https://github.com/approov/documentation_backend-code-snippets/blob/master/golang/src/backend-integration-impact/hello-server-token-binding-protected.go

package main

import (
    "os"
    "fmt"
    "encoding/json"
    "encoding/base64"
    "log"
    "net/http"
    jwt "github.com/dgrijalva/jwt-go"
    "crypto/sha256"
)

var base64Secret = os.Getenv("APPROOV_BASE64_SECRET")

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

func verifyApproovTokenBinding(token *jwt.Token, request *http.Request) (jwt.Claims, error) {
    claims := token.Claims
    token_binding_payload, has_pay_key := claims.(jwt.MapClaims)["pay"]

    if ! has_pay_key {
        // We don't return an error, because the Approov fail-over server doesn't provide the `pay` key.
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

        claims, err := verifyApproovTokenBinding(token, request)

        if err != nil {
            errorResponse(response, http.StatusUnauthorized, err.Error())
            return
        }

        // Use the returned claims to perform any additional checks.
        _ = claims

        handler(response, request)
    })
}

func main() {
    http.Handle("/", makeApproovCheckerHandler(helloHandler))
    log.Println("Server listening on http://localhost:8002")
    err := http.ListenAndServe(":8002", nil)
    if err != nil {
        log.Fatal("Server Error: " + err.Error())
    }
}
