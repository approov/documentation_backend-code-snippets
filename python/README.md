# PYTHON CODE SNIPPETS

Python code snippets we will use in the [Approov 2 documentation](https://approov.io/docs/).


## SETUP

#### Clone

```
git clone git@github.com:approov/documentation_backend-code-snippets.git && cd documentation_backend-code-snippets/python
```

#### Install Dependencies

From the root of the Python project run:

```bash
pip install --user -r requirements.txt
```

### Environment

Inside each `src` folder exists a `.env.example` file that you will need to copy to `.env`.

The `APPROOV_BASE64_SECRET` in the `.env.example` file was generated with `openssl rand -base64 64 | tr -d '\n'; echo`,
and the JWT tokens used in the [Postman Collection](./../api.postman_collection.json) of this project were created in the
https://jwt.io site. To create your own tokens you can use [this example](https://jwt.io/#debugger-io?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MDg2ODMyMDUuODkxOTEyfQ.c8I4KNndbThAQ7zlgX4_QDtcxCrD9cff1elaCJe9p9U),
that doesn't contain the Approov token binding, or you can use the [token binding](https://jwt.io/#debugger-io?token=eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjQ3MTgwMTgyMjQuNzgwMzY4LCJwYXkiOiJWUUZGUEpaNjgyYU90eFJNanowa3RDSG15V2VFRWVTTXZYaDF1RDhKM3ZrPSJ9.N-KwuLeUt9s6TDibhX32AIkhobCYVh5-brVESqUxdBk) example.
Please don't forget to provide the secret to sign the JWT token, that is whatever you have defined in `APPROOV_BASE64_SECRET`.

> **NOTE:** For production usage the secret is always retrieved with the Approov CLI tool, that can be also used to
            generate valid tokens for testing purposes. Check the Approov CLI tool docs [here](https://approov.io/docs/v2.1/approov-cli-tool-reference/#token-commands).


## HOW TO RUN THE CODE SNIPPETS

To run any of the example server just go inside its folder and run `flask run`, if inside docker you need to run `flask run -h 0.0.0.0` in order to expose the server outside the container network.

```
cd src/example-api-integration && flask run
```

or

```
cd src/backend-integration-impact && flask run
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
HTTP/1.0 200 OK
Content-Type: application/json
Content-Length: 32
Server: Werkzeug/0.16.0 Python/3.7.4
Date: Tue, 24 Sep 2019 13:26:16 GMT

{
  "message": "Hello World!"
}
```
