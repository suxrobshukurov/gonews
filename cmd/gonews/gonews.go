package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/suxrobshukurov/gonews/pkg/api"
	"github.com/suxrobshukurov/gonews/pkg/rss"
	"github.com/suxrobshukurov/gonews/pkg/storage"
	"github.com/suxrobshukurov/gonews/pkg/storage/memdb"
	"github.com/suxrobshukurov/gonews/pkg/storage/postgres"
)

type gonewsConfig struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// init requirements
	var srv server
	var err error
	srv.db, err = postgres.New()
	if err != nil {
		srv.db, _ = memdb.New()
		log.Printf("failed to create a postgress database: %v\n memdb was launched instead", err)
	}
	srv.api = api.New(srv.db)

	// read config
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("failed to read config: ", err)
	}

	var config gonewsConfig
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal("failed to unmarshal config: ", err)
	}

	// parse urls in goroutine
	chPosts := make(chan []storage.Post)
	chErros := make(chan error)
	for _, url := range config.URLS {
		go parseUrl(url, chPosts, chErros, config.Period)
		log.Println("started parsing url: ", url)
	}

	// read channels and send to db
	go func() {
		for posts := range chPosts {
			err = srv.db.AddPosts(posts)
			if err != nil {
				log.Printf("failed to add posts: %v", err)
			} else {
				log.Printf("added %d posts", len(posts))
			}
		}
	}()

	// read errors and log them
	go func() {
		for err := range chErros {
			log.Println("failed to parse url: ", err)
		}
	}()

	log.Println("started server: http://localhost:80")
	err = http.ListenAndServe(":80", srv.api.Router())
	if err != nil {
		log.Fatal("failed to start server: ", err)
	}

}

func parseUrl(url string, chPosts chan<- []storage.Post, chErrors chan<- error, period int) {
	for {
		posts, err := rss.ParseRSS(url)
		if err != nil {
			chErrors <- fmt.Errorf("failed to parse %s: %w", url, err)
			continue
		}
		chPosts <- posts
		time.Sleep(time.Minute * time.Duration(period))
	}
}
