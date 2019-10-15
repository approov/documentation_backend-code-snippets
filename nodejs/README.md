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

#### Environment

```
cp .env.example .env
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
X-Powered-By: Express
Content-Type: application/json; charset=utf-8
Content-Length: 39
ETag: W/"27-u2HmDGig9dI3I4Ws5u7Bi/8Jh7o"
Date: Fri, 02 Aug 2019 13:06:55 GMT
Connection: keep-alive

{"message":"Hello Approov token binding protected server!"}
```
