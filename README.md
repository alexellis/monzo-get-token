## monzo-get-token

Golang example project for requesting / retrieving a Monzo API token

This code uses the OAuth 2.0 code grant.

## Pre-reqs:

* Create an App on [Monzo developer portal](https://docs.monzo.com/)

* Install Golang 1.10 or newer

* Clone this GitHub repo

## 1. Build the code and insert your tokens

```
export client_id=""
export client_secret=""

export redirect_uri="http://localhost:8088/oauth2/callback"
export port=8088

go build && ./monzo-get-token
```

## 2. Start the OAuth 2.0 flow

Navigate to http://localhost:8088

## 3. Authorize the app and retrieve the token

The token will be displayed in the console log and in your browser.

You can now use it to access the [Monzo API](https://docs.monzo.com/).

## 4. Operationalize

Now you have the token, pay attention to the `expires_in` value - you'll need to use the `refresh_token` to renew your access token when it expires.

License: MIT

Copyright Alex Ellis 2018