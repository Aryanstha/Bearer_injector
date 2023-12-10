package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// User structure
type User struct {
	ClientID string
	Salt     string
}

var users = map[string]User{
	"1": {"0388941f", "vinegar"},
	"2": {"4be75c87", "tabasco"},
	"3": {"12345678", "hendos"},
}

func genHash(path, clientID, salt, body string) string {
	hash := sha256.New()
	hash.Write([]byte(path + clientID + body + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("user-id")
		log.Println("User-id:", userID)

		if userID == "-1" {
			return
		}

		var user User
		var ok bool

		if user, ok = users[userID]; !ok {
			log.Println("Unknown user requested, using user 1")
			user = users["1"]
		} else {
			log.Println("Using user", userID)
		}

		urlParsed, err := url.Parse(r.URL.String())
		if err != nil {
			log.Println("Error parsing URL:", err)
			return
		}

		path := urlParsed.Path
		log.Println("URL:", r.URL.String())
		log.Println("Path:", path)
		log.Println("Client ID:", user.ClientID)
		log.Println("Salt:", user.Salt)

		body := r.FormValue("body") // Assuming the body is passed as a form value
		log.Println("Body:", body)

		token := genHash(path, user.ClientID, user.Salt, body)
		log.Println("Token:", token)

		if r.Header.Get("bearer") == "" {
			r.Header.Set("bearer", token)
		} else {
			log.Println("Bearer token already passed, not modifying it")
		}

		if r.Header.Get("client-id") == "" {
			r.Header.Set("client-id", user.ClientID)
		} else {
			log.Println("Client ID already passed, not modifying it")
		}

		w.Write([]byte("OK"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
