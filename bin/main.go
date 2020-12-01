package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/dereulenspiegel/go-brewchild"
)

func main() {

	downloadBatches()

	imageLinks := scrapeUntappdImages()

	generateContent(imageLinks)
}

func generateContent(imageLinks map[string][]string) {
	for _, b := range str.Beers {
		beerContentFolder, err := ensureBeerContentFolder(b.UntappdID, b.Name)
		if err != nil {
			log.Fatalf("Failed to create content folder for beer %s: %s", b.UntappdID, err)
		}
		if err := ensureBeerContent(beerContentFolder, b); err != nil {
			log.Fatalf("Failed to create content in folder %s: %s", beerContentFolder, err)
		}
		if err := downloadImages(imageLinks[b.UntappdID], beerContentFolder); err != nil {
			log.Printf("[ERROR] %s", err)
		}

		for _, bt := range b.Batches {
			ensureBatchPostData(bt)
		}
	}

	for _, bt := range str.UnlinkedBatches {
		ensureBatchPostData(bt)
	}
}

func downloadImages(links []string, targetFolder string) error {
	for _, l := range links {
		if err := downloadImage(l, targetFolder); err != nil {
			log.Printf("[ERROR] Failed to download %s to %s", l, targetFolder)
			continue
		}
	}
	return nil
}

func downloadImage(link, targetFolder string) error {
	log.Printf("Downloading image %s to %s", link, targetFolder)
	imageName := path.Base(link)
	if imageName == "badge-beer-default.png" {
		log.Printf("Skipping default beer batch")
		return nil
	}
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("Failed to download %s: %w", link, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Received invalid status code %d for link %s", resp.StatusCode, link)
	}
	imagePath := filepath.Join(targetFolder, imageName)
	imageFile, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("Failed to create image file %s: %w", imagePath, err)
	}
	defer imageFile.Close()
	if _, err := io.Copy(imageFile, resp.Body); err != nil {
		return fmt.Errorf("Failed to download %s to %s: %w", link, imagePath, err)
	}
	return nil
}

func downloadBatches() {
	bfClient, err := brewchild.New(brewfatherUserID, brewfatherAPIKey)
	if err != nil {
		log.Fatalf("Failed to create brewfather client: %s", err)
	}

	for _, state := range []string{"Completed", "Planning", "Brewing", "Fermenting", "Conditioning", "Archived"} {
		log.Printf("Downloading %s batches", state)
		batches, err := bfClient.Batches(brewchild.Status(state), brewchild.Complete(true), brewchild.Limit(100))
		if err != nil {
			log.Fatalf("Failed to retrieve batches from brewfather: %s", err)
		}
		log.Printf("Got %d %s batches", len(batches), state)
		myBatches := exportBatches(batches, state)
		for _, m := range myBatches {
			log.Printf("Adding batch number %d with untappdID %s", m.Number, m.UntappdID)
			if err := str.AddBatch(m); err != nil {
				log.Fatalf("Failed to add batch: %s", err)
			}
		}
	}
}

func scrapeUntappdImages() map[string][]string {
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

	linkColl := map[string][]string{}

	for _, item := range feed.Channel.Items {
		beerID, imageLink := parseCheckin(item.Link)
		if beerID == "" {
			log.Printf("Didn't receive usabel data for checkin %s", item.Link)
			continue
		}
		links := linkColl[beerID]
		links = append(links, imageLink)
		linkColl[beerID] = links
		if _, exists := str.Beers[beerID]; !exists {
			log.Printf("Found not existing beer %s in untappd feed, scraping it", beerID)
			br, err := scrapeBeerDetails(beerID)
			if err != nil {
				log.Printf("[ERROR] Failed to scrape beer infos for untappd id %s: %s", beerID, err)
				continue
			}
			if err := str.AddBeer(br); err != nil {
				log.Printf("[ERROR] Failed to add beer to store: %s", err)
				continue
			}
		}
	}
	return linkColl
}
