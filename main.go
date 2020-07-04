package main

import (
	"bufio"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Link Fuck
type Link struct {
	Name string
	URL  string
}

var urlLinks = []Link{
	Link{Name: "about", URL: "/about"},
	Link{Name: "contact", URL: "/contact"},
	Link{Name: "sites", URL: "/sites"},
	Link{Name: "privacy", URL: "/privacy"},
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getQuotes() []string {
	// Create a new list
	quotes := make([]string, 0)
	// Create a new variable called files, which loads assets/quotes
	file, err := os.Open("assets/quotes")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		quotes = append(quotes, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return quotes
}

func randQuote(quotes []string) string {
	randIntQuotes := rand.Intn(len(quotes))
	quote := quotes[randIntQuotes]
	return quote
}

func main() {
	quotes := getQuotes()
	// Initialize a new instance of gin.
	// := is shorthand for var router = gin.Default()
	router := gin.Default()

	// Glob all files within assets/templates
	router.LoadHTMLGlob("assets/templates/*")
	// Declare folder for static assets
	// maps assets/static to /src
	router.Static("/css", "static/css")
	router.Static("/js", "static/js")
	router.Static("/static", "static/meta")
	router.Static("/img", "static/img")
	// Create a path for the home
	router.GET("/", func(c *gin.Context) {
		quote := randQuote(quotes)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"name":       "home",
			"headerText": "oxide.one",
			"quote":      quote,
			"urlLinks":   urlLinks,
		})
	})
	// About Page
	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl", gin.H{
			"name":       "about",
			"headerText": "About",
			"urlLinks":   urlLinks,
		})
	})
	// Contact page
	router.GET("/contact", func(c *gin.Context) {
		c.HTML(http.StatusOK, "contact.tmpl", gin.H{
			"name":       "contact",
			"headerText": "Contact Me",
			"urlLinks":   urlLinks,
		})
	})
	// Sites Page
	router.GET("/sites", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sites.tmpl", gin.H{
			"name":       "sites",
			"headerText": "Sites",
			"urlLinks":   urlLinks,
		})
	})
	// Privacy Page
	router.GET("/privacy", func(c *gin.Context) {
		c.HTML(http.StatusOK, "privacy.tmpl", gin.H{
			"headerText": "Privacy",
			"urlLinks":   urlLinks,
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func init() {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().Unix())
}
