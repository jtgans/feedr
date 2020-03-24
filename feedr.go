package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/feeds"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var Version string = "0.1"

var MaxItemCount = flag.Int("max-item-count", 100, "Maximum number of items to keep in RAM")
var FeedTitle = flag.String("title", "Feedr Newsfeed", "The title of the RSS feed")
var FeedLink = flag.String("link", "", "Link to the RSS feed")
var AuthorName = flag.String("author-name", "", "Name of the feed author")
var AuthorEmail = flag.String("author-email", "", "Email address of the feed author")
var Description = flag.String("description", "", "Description of the feed")

type FeedItem struct {
	Title       string
	Link        string
	Source      string
	AuthorName  string
	AuthorEmail string
	Description string
	Content     string
}

func handleGetRequest(f *feeds.Feed, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/rss+xml")
	err := f.WriteRss(w)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("E GET %v Failure attempting to serialize feed: %v\n", r.RemoteAddr, err)
	}

	log.Printf("I GET %v 200\n", r.RemoteAddr)
}

func handlePostRequest(f *feeds.Feed, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("E POST %v Error reading body: %v\n", r.RemoteAddr, err)
		return
	}

	var result FeedItem
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Item deserialization error", http.StatusBadRequest)
		log.Printf("E POST %v Error parsing request: %v. Request was: %v\n", r.RemoteAddr, err, string(body))
		return
	}

	if (result.Title == "") && (result.Description == "") {
		http.Error(w, "Missing item title or description", http.StatusBadRequest)
		log.Printf("E POST %v New item request missing title or description. Request was: %v\n", r.RemoteAddr, string(body))
		return
	}

	id := feeds.NewUUID().String()
	i := &feeds.Item{
		Title:       result.Title,
		Link:        &feeds.Link{Href: result.Link},
		Source:      &feeds.Link{Href: f.Link.Href},
		Author:      &feeds.Author{Name: result.AuthorName, Email: result.AuthorEmail},
		Description: result.Description,
		Content:     result.Content,
		Created:     time.Now(),
		Id:          id,
	}
	f.Add(i)
	f.Updated = time.Now()
	f.Created = time.Now()

	if len(f.Items) > *MaxItemCount {
		f.Items = f.Items[1:]
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v\n", id)
	log.Printf("I POST %v 201 Created new item %v\n", r.RemoteAddr, id)
}

func main() {
	log.Printf("I Feedr v%v starting...\n", Version)
	flag.Parse()

	if *FeedLink == "" {
		log.Fatal("A feed link must be provided.")
	}
	if *FeedTitle == "" {
		log.Fatal("A feed title must be provided.")
	}
	if *Description == "" {
		log.Fatal("A description must be provided.")
	}

	f := &feeds.Feed{
		Title:       *FeedTitle,
		Link:        &feeds.Link{Href: *FeedLink},
		Description: *Description,
	}
	if (*AuthorName != "") || (*AuthorEmail != "") {
		(*f).Author = &feeds.Author{Name: *AuthorName, Email: *AuthorEmail}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.RequestURI != "/" {
				log.Printf("I %v %v %v 404", r.Method, r.RequestURI, r.RemoteAddr)
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				handleGetRequest(f, w, r)
			}
		} else if r.Method == "POST" {
			handlePostRequest(f, w, r)
		} else {
			log.Printf("I %v %v %v 501", r.Method, r.RequestURI, r.RemoteAddr)
			http.Error(w, "Method not supported", http.StatusNotImplemented)
		}
	})

	log.Println("I Serving on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
