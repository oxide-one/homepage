package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Link Structure
type Link struct {
	Name string
	URL  string
}

type DataResponse struct {
	Monitors []struct {
		FriendlyName       string `json:"friendly_name"`
		URL                string `json:"url"`
		Status             int    `json:"status"`
		CustomUptimeRanges string `json:"custom_uptime_ranges"`
	} `json:"monitors"`
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

var lastChecked = time.Now().Add(-10 * time.Minute)
var lastStatus DataResponse

func getStatus(apiKey string) DataResponse {
	// Figure out if the time difference is over 10 minutes
	timeSince := time.Since(lastChecked)
	if timeSince <= time.Minute*10 {
		return lastStatus
	}
	// Figure out what date is today
	now := time.Now()

	// Get today's date in Epoch time
	nowSecs := now.Unix()
	// Figure out what date was yesterday
	past := now.AddDate(0, 0, -7)
	// Get 7 days ago in Epoch time
	pastSecs := past.Unix()

	// Set the URL for uptimerobot
	url := "https://api.uptimerobot.com/v2/getMonitors"

	payload := strings.NewReader(fmt.Sprintf("api_key=%s&format=json&custom_uptime_ranges=%d_%d", apiKey, pastSecs, nowSecs))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	response, err := http.DefaultClient.Do(req)
	statusresponse := DataResponse{}
	if err != nil {
		fmt.Printf("%s", err)
		//os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		json.Unmarshal(contents, &statusresponse)
	}
	// Set the cache status
	lastChecked = time.Now()
	lastStatus = statusresponse
	return statusresponse
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

func getAPIKey() string {
	// Checking that the env var "API_KEY" is present.
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		log.Fatal("API_KEY IS NOT SET")
	}
	return apiKey
}

func main() {
	quotes := getQuotes()
	// Initialize a new instance of gin.
	// := is shorthand for var router = gin.Default()
	router := gin.Default()

	// Grab environment variables
	apiKey := getAPIKey()

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
		status := getStatus(apiKey)
		c.HTML(http.StatusOK, "sites.tmpl", gin.H{
			"status":     status.Monitors,
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
