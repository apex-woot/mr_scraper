package scraper

import (
	"context"
	"log"

	"github.com/apexwoot/mr_scraper/config"
	"github.com/apexwoot/mr_scraper/internal/app"
)

func Scrape() {
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading env: %s", err)
	}
	dbConn, err := config.NewConnection(config.ComposeDBConnectionString())
	if err != nil {
		log.Fatalf("Error while creating DB connection: %v", err)
	}
	defer dbConn.Disconnect(context.TODO())

	newestPost, err := dbConn.GetNewestPost()
	if err != nil {
		log.Fatalf("Error retrieving newest post: %v", err)
	}
	newestCheckedPostTimestamp := app.GetNewestPostTimestampOrDefault(newestPost)

	doc, err := app.FetchHTML("https://www.reddit.com/r/mechmarket/search/?q=gmk&sort=new")
	if err != nil {
		log.Fatalf("Error fetching HTML: %s", err)
	}

	element := app.FindTagByName(doc, "reddit-feed")
	posts, err := app.GetChildren(element, newestCheckedPostTimestamp)
	if err != nil {
		log.Fatalf("Error while retrieving posts: %v", err)
	}

	log.Printf("Found %d new posts. Saving...", len(posts))

	_, err = dbConn.SavePosts(posts)
	if err != nil {
		log.Fatalf("Error saving posts: %v", err)
	}
	log.Printf("Saved %d new posts.", len(posts))
}
