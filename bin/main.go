package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"sync"

	"github.com/dereulenspiegel/go-brewchild"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		downloadBatches()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		scrapeUntappdImages()
	}()
	wg.Wait()
}

func downloadBatches() {
	bfClient, err := brewchild.New(brewfatherUserID, brewfatherAPIKey)
	if err != nil {
		log.Fatalf("Failed to create brewfather client: %w", err)
	}

	for _, state := range []string{"Completed", "Planning", "Brewing", "Fermenting", "Conditioning", "Archived"} {
		exportBatches(bfClient, state)
	}
}

func scrapeUntappdImages() {
	resp, err := http.Get(breweryFeed)
	if err != nil {
		log.Fatalf("Failed to get brewery feed: %s", err)
	}
	defer resp.Body.Close()

	feed := feedXML{}
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		log.Fatalf("Failed to parse brewery feed: %s", err)
	}

	wg := &sync.WaitGroup{}
	for _, item := range feed.Channel.Items {
		wg.Add(1)
		go func(checkinURL string) {
			defer wg.Done()
			parseCheckin(checkinURL)
		}(item.Link)
	}
	wg.Wait()
}
