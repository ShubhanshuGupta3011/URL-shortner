package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data: ", data)
	hash := hex.EncodeToString(data)
	fmt.Println("EncodeToString: ", hash)
	fmt.Println("final string: ", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	// var data struct {
	// 	URL string `json:"url"`
	// }
	// err := json.NewDecoder(r.Body).Decode(&data)
	// if err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	URL := r.URL.Query().Get("url")

	shortURL_ := createURL(URL)
	// fmt.Fprintf(w, shortURL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	OriginalURL := "https://leetcode.com/u/SHUBHANSHU_GUPTA-156/"
	createURL(OriginalURL)
	// ShortURL = 287d2e56
	// http.//localhost:3000/redirect/287d2e56
	// working correctly but not showing in postman
	http.HandleFunc("/", corsMiddleware(RootPageURL))
	http.HandleFunc("/shorten", corsMiddleware(ShortURLHandler))
	http.HandleFunc("/redirect/", corsMiddleware(redirectURLHandler))
	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe("localhost:3000", nil)
	if err != nil {
		fmt.Println("Error on starting server:", err)
	}
}
