package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type URLShortener struct {
	urls map[string]string
}

var (
	// Create a logger to write to a file
	logger *log.Logger
)

// Generates the HTML for the URL shortener forms
func generateForm(originalURL, shortenedURL string) string {
	return fmt.Sprintf(`<div style="text-align:center">
	<h2>URL Shortener</h2>
	<p>Original URL: %s</p>
	<p>Shortened URL: <a href="%s">%s</a></p>
	<form method="post" action="/shortener">
		<input type="text" name="url" placeholder="Enter a URL" required>
		<input type="submit" value="Shorten">
	</form>
	</div>`, originalURL, shortenedURL, shortenedURL)
}

// HandleShortener handles both GET and POST requests for URL shortening
func (us *URLShortener) HandleShortener(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
	if r.Method == http.MethodGet {
		// Display the form
		w.Header().Set("Content-Type", "text/html")
		formHTML := generateForm("", "")
		fmt.Fprint(w, formHTML)
		logger.Println("Displayed form")
		return
	}

	if r.Method == http.MethodPost {
		// Handle URL shortening
		originalURL := r.FormValue("url")
		if originalURL == "" {
			http.Error(w, "URL parameter is missing", http.StatusBadRequest)
			logger.Println("Error: URL parameter is missing")
			return
		}

		shortKey := generateShortKey()
		us.urls[shortKey] = originalURL

		shortenedURL := fmt.Sprintf("http://localhost:3000/short/%s", shortKey)

		w.Header().Set("Content-Type", "text/html")
		formHTML := generateForm(originalURL, shortenedURL)
		fmt.Fprint(w, formHTML)
		logger.Printf("Shortened URL: %s to %s", originalURL, shortenedURL)
		return
	}

	// Handle other methods
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	logger.Println("Error: Invalid request method")
}

// HandleRedirect handles URL redirection
func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		logger.Println("Error: Shortened key is missing")
		return
	}

	originalURL, found := us.urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		logger.Printf("Error: Shortened key %s not found", shortKey)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	logger.Printf("Redirected %s to %s", shortKey, originalURL)
}

// GenerateShortKey generates a random key for the URL
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

// Middleware to log requests and errors
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Open the log file
	file, err := os.OpenFile("errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	// Set up the logger
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println("Starting URL shortener...")

	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	// Register handlers with logging middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/shortener", shortener.HandleShortener)
	mux.HandleFunc("/short/", shortener.HandleRedirect)

	// Use the logging middleware
	http.ListenAndServe(":3000", loggingMiddleware(mux))
}
