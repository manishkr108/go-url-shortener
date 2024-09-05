# go-url-shortener
A simple URL shortener built with Go. This application allows users to shorten long URLs and redirects them to the original URL using a short key.



Here’s a descriptive README for your URL shortener application:

URL Shortener
A simple URL shortener built with Go. This application allows users to shorten long URLs and redirects them to the original URL using a short key.

**Features**
**Shorten URLs**: Input a long URL and get a shortened version that redirects to the original URL.<br>
**Generate Unique Keys**: Automatically generates unique short keys for each URL.<br>
**Redirect to Original URL**: Redirects users to the original URL when they visit the shortened URL.<br>
**Installation**<br>
Clone the Repository:


**Copy code**
git clone https://github.com/manishkr108/go-url-shortener.git <br>
cd go-url-shortener<br>
Build the Application:<br>

sh
Copy code
go build -o go-url-shortener<br>
Run the Application:<br>

sh
Copy code
./go-url-shortener<br>
Access the Application: Open your web browser and navigate to http://localhost:3000/shortener to use the URL shortener.

Usage
Shorten a URL
Open http://localhost:3000/shortener in your web browser.
Enter the URL you want to shorten in the provided form and click "Shorten."
The page will display the shortened URL which you can use to access the original URL.
Redirect to Original URL
Visit the shortened URL (e.g., http://localhost:3000/short/{shortKey}).
You will be redirected to the original URL associated with the short key.
Code Explanation
URLShortener Struct: Manages a map of short keys to original URLs.

HandleShortener Method: Handles both GET and POST requests:

GET Request: Displays the URL shortening form.
POST Request: Processes the URL shortening form submission, generates a short key, and displays the shortened URL.
HandleRedirect Method: Handles redirects from shortened URLs to the original URL based on the short key.

generateShortKey Function: Generates a random 6-character string to use as the short key.
