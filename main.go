package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("/", makeHandler())

	router.HandleFunc("/oauth2/callback", makeCallbackHandler())

	router.HandleFunc("/healthz/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK."))
	})

	timeout := time.Second * 10
	port := 8080
	if v, exists := os.LookupEnv("port"); exists {
		val, _ := strconv.Atoi(v)
		port = val
	}

	log.Printf("Using port: %d\n", port)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        router,
		ReadTimeout:    timeout,
		WriteTimeout:   timeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func makeCallbackHandler() func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		code := q.Get("code")

		if len(code) > 0 {
			form := url.Values{}
			form.Set("grant_type", "authorization_code")
			form.Set("client_id", os.Getenv("client_id"))
			form.Set("client_secret", os.Getenv("client_secret"))
			form.Set("redirect_uri", os.Getenv("redirect_uri"))
			form.Set("code", code)

			req, reqErr := http.NewRequest(http.MethodPost, "https://api.monzo.com/oauth2/token", strings.NewReader(form.Encode()))
			if reqErr != nil {
				panic(reqErr)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			res, resErr := http.DefaultClient.Do(req)

			if resErr != nil {
				panic(resErr)
			}

			log.Printf("Status: %d", res.StatusCode)
			body, _ := ioutil.ReadAll(res.Body)
			if res.Body != nil {
				defer res.Body.Close()
			}
			log.Printf("Body: %s", body)
			w.Write(body)
		}
	}
}

func makeHandler() func(w http.ResponseWriter, r *http.Request) {
	// https://auth.monzo.com/?
	// client_id=$client_id&
	// redirect_uri=$redirect_uri&
	// response_type=code&
	// state=$state_token

	return func(w http.ResponseWriter, r *http.Request) {
		authURL, _ := url.Parse("https://auth.monzo.com/")
		q := authURL.Query()
		q.Set("client_id", os.Getenv("client_id"))
		q.Set("redirect_uri", os.Getenv("redirect_uri"))
		q.Set("response_type", "code")
		q.Set("state", fmt.Sprintf("%d", time.Now().Unix()))

		authURL.RawQuery = q.Encode()

		http.Redirect(w, r, authURL.String(), http.StatusTemporaryRedirect)
	}
}
