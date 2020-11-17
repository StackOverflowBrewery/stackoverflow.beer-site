package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var (
	breweryFeed = "https://untappd.com/rss/brewery/462497"

	untappdImageFolder = "../static/images/untappd"
)

type feedXML struct {
	Channel struct {
		Items []feedItem `xml:"item"`
	} `xml:"channel"`
}

type feedItem struct {
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func main() {
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

type checkinPage struct {
	Metas []pageMeta `xml:"meta"`
}

type pageMeta struct {
	Property string `xml:"property,attr"`
	Content  string `xml:"content,attr"`
}

func parseCheckin(checkinURL string) {
	resp, err := http.Get(checkinURL)
	if err != nil {
		log.Fatalf("Failed to retrieve checkin page: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse checkin page")
	}

	var imageLink string
	var description string

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("property")
		if exists && val == "og:image" {
			contentVal, exists := s.Attr("content")
			if exists {
				imageLink = contentVal
			}
		}

		if exists && val == "og:description" {
			contentVal, exists := s.Attr("content")
			if exists {
				description = contentVal
			}
		}

	})
	if imageLink != "" && !strings.HasSuffix(imageLink, ".png") {
		resp, err := http.Get(imageLink)
		if err != nil {
			log.Fatalf("Failed to download checkin image")
		}
		defer resp.Body.Close()
		imageName := path.Base(imageLink[7:])
		outPath := filepath.Join(untappdImageFolder, imageName)
		outFile, err := os.Create(outPath)
		if err != nil {
			log.Fatalf("Failed to create output image: %s", err)
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			log.Fatalf("Failed to write image to disk: %s", err)
		}
	}
}
