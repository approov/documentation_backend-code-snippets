# GOLANG CODE SNIPPETS

Golang code snippets to be used in the [Approov 2 documentation](https://approov.io/docs/).

## SETUP

#### Clone

```
git clone git@github.com:approov/documentation_backend-code-snippets.git && cd documentation_backend-code-snippets/golang
```

### Environment

Along side each `go.mod` file exists a `.env.example` file that contains the `APPROOV_BASE64_SECRET` for the JWT tokens
used in the [Postman Collection](./../api.postman_collection.json) that we have in the root of this project.

The `APPROOV_BASE64_SECRET` needs to be present in your environment before you can run the go server, and in Linux you
can do it like:

```
set -o allexport && source .env.example && set +o allexport
```

The `APPROOV_BASE64_SECRET` in the `.env.example` file was generated with `openssl rand -base64 64 | tr -d '\n'; echo`,
and the JWT tokens used in the [Postman Collection](./../api.postman_collection.json) of this project were created in the
https://jwt.io site. To create your own tokens you can use [this example](https://jwt.io/#debugger-io?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MDg2ODMyMDUuODkxOTEyfQ.c8I4KNndbThAQ7zlgX4_QDtcxCrD9cff1elaCJe9p9U),
that doesn't contain the Approov token binding, or you can use the [token binding](https://jwt.io/#debugger-io?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MTgwMTgyMjQuNzgwMzY4LCJwYXkiOiJWUUZGUEpaNjgyYU90eFJNanowa3RDSG15V2VFRWVTTXZYaDF1RDhKM3ZrPSJ9.N-KwuLeUt9s6TDibhX32AIkhobCYVh5-brVESqUxdBk) example.
Please don't forget to provide the secret to sign the JWT token, that is whatever you have defined in `APPROOV_BASE64_SECRET`.

> **NOTE:** For production usage the secret is always retrieved with the Approov CLI tool, that can be also used to
            generate valid tokens for testing purposes. Check the Approov CLI tool docs [here](https://approov.io/docs/v2.1/approov-cli-tool-reference/#token-commands).


## HOW TO RUN THE CODE SNIPPETS

To run the several code snippets and test them we can start one of the following servers:

```
cd src/example-api-integration/unprotected && go run hello-server-unprotected.go
```
or
```
cd src/example-api-integration/protected && go run hello-server-protected.go
```
or
```
cd src/backend-integration-impact && go run hello-server-token-binding-protected.go
```

To interact with the server just use the Postman collection that you can download from [here](./../api.postman_collection.json).


#### Approov Token Binding Example

When the request is made with a valid `Approov-Token` that also contains a value in the key `pay` that matches the `Authorization` header, the request will be considered to come from a genuine mobile app, and to simulate it we can issue a request form a tool like Postman or Curl.

The request from Postman:

![Valid Approov Token Binding Request Example](./../.assets/img/postman-valid-approov-token-binding.png)

But you can test it also with CURL:

```
curl -iX GET \
  http://localhost:8002/ \
  -H 'Approov-Token: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MTgwMTgyMjQuNzgwMzY4LCJwYXkiOiJWUUZGUEpaNjgyYU90eFJNanowa3RDSG15V2VFRWVTTXZYaDF1RDhKM3ZrPSJ9.N-KwuLeUt9s6TDibhX32AIkhobCYVh5-brVESqUxdBk' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c' \
  -H 'cache-control: no-cache'
```

That will receive this response from the server:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 05 Aug 2019 21:28:49 GMT
Content-Length: 28

{"message":"Hello, World!"}
```
