package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/FilipStam97/rssator/internal/database"
)

// we have a pointer to the db, number of go routines that will do the scraping and the time between the scrapes
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	//empty for start and middle since wwe want the first iteration to fire immediately, so the loop will run every tick (if we used range it would wait the first tick)
	for ; ; <-ticker.C { //infinite for loop
		feeds, err := db.GetNextFeedToFetch(
			context.Background(), // context.Background is a global context, what you use if you don't have access to scoped context
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue //we are continuing because we want our scraper to always be running if we return the function will exit
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1) //adds one to the wait group so when we get to the wg.Wait it will wait for x amount of distinct calls to wg.Done , done decrements the wg counter
			//this allows us to spawn eg. 30 different go routines at the same time to scrape different feeds
			//we are iterating over the all feeds that we want to fetch via goroutines and adding them to the wait group

			go scrapeFeed(db, feed, wg) //spawns a new goroutine in the context of the wait group
		}
		wg.Wait() //when all are done this line will stop blocking
	}
}

func scrapeFeed(db *database.Queries, feed database.Feed, wg *sync.WaitGroup) {
	defer wg.Done() //calling done when we scrape the feed

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	//logging to console for now instead of storing to db
	for _, item := range rssFeed.Chanel.Item {
		log.Println("Found post", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Chanel.Item))

}
