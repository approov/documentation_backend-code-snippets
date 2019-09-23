# NODEJS CODE SNIPPETS

NodeJS code snippets we will use in the [Approov 2 documentation](https://approov.io/docs/).

## SETUP

#### Clone

```
git clone git@github.com:approov/documentation_backend-code-snippets.git && cd documentation_backend-code-snippets/nodejs
```

#### Install Dependencies

```
npm install
```

## HOW TO RUN THE CODE SNIPPETS

To run the several code snippets and test them we can start one of the following servers:

```
npm run unprotected-server
```
or
```
npm run protected-server
```
or
```
npm run token-binding-protected-server
```

To interact with the server just use the Postman collection that you can download from [here](./../api.postman_collection.json).


#### Approov Token Binding Example

This example is for using the request in Postman for an `Approov-Token` with a `pay` key that matches the `Authorization` header:

```
$ npm run token-binding-protected-server

> approov2-code-snippets@1.0.0 token-binding-protected-server /home/node/workspace
> node src/backend-integration-impact/hello-server-token-binding-protected.js

Server listening on http://localhost:8002
VALID APPROOV TOKEN
VALID APPROOV TOKEN BINDING.
```

But you can test it also with CURL:

```
curl -iX GET \
  http://localhost:8002/ \
  -H 'Approov-Token: eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MTgwMTgyMjQuNzgwMzY4LCJwYXkiOiJWUUZGUEpaNjgyYU90eFJNanowa3RDSG15V2VFRWVTTXZYaDF1RDhKM3ZrPSJ9.N-KwuLeUt9s6TDibhX32AIkhobCYVh5-brVESqUxdBk' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c' \
  -H 'cache-control: no-cache'
```

That will output:

```
HTTP/1.1 200 OK
X-Powered-By: Express
Content-Type: application/json; charset=utf-8
Content-Length: 39
ETag: W/"27-u2HmDGig9dI3I4Ws5u7Bi/8Jh7o"
Date: Fri, 02 Aug 2019 13:06:55 GMT
Connection: keep-alive

{"message":"Hello Approov token binding protected server!"}
```
