package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type URLShortner struct {
	urls map[string]string
}

// ! Implement URL Redirection
func (us *URLShortner) HandleShortener(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(w, "URL Parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:3000/short/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
	<h2> URL Shortener</h2>
	<p> Original URL: %s</p>
	<p> Shortened URL:<a href="%s">%s </a></p>
	<form method="post" action="/shortener">
	<input type="text" name="url" placeholder="Enter a URL">
	<input type="submit" value="Shorten">
	</form>
	`, originalURL, shortenedURL, shortenedURL)
	fmt.Fprint(w, responseHTML)
}

//!handel URL Redirect

func (us *URLShortner) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

// ! Generate Short Key

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]

	}

	return string(shortKey)
}

func main() {
	fmt.Println("Statring URL shortner....")

	shortener := &URLShortner{
		urls: make(map[string]string),
	}
	http.HandleFunc("/shortener", shortener.HandleShortener)
	http.HandleFunc("/short/", shortener.HandleRedirect)

	fmt.Println("URL Shortener is running on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
