package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func genHash(path, clientID, salt, body string) string {
	hash := sha256.New()
	hash.Write([]byte(path + clientID + body + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

func request(w http.ResponseWriter, r *http.Request) {
	clientID := "0388941f"
	salt := "vinegar"

	urlParsed, err := url.Parse(r.URL.String())
	if err != nil {
		log.Println("Error parsing URL:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	path := urlParsed.Path
	body := r.FormValue("body") // Assuming the body is passed as a form value

	log.Println("URL:", r.URL.String())
	log.Println("Path:", path)
	log.Println("Client ID:", clientID)
	log.Println("Body:", body)

	token := genHash(path, clientID, salt, body)
	log.Println("Token:", token)

	w.Header().Set("bearer", token)
	w.Header().Set("client-id", clientID)

	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", request)

	port := 8080
	log.Printf("Server listening on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
