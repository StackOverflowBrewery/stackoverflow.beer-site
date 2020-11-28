package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	breweryFeed = "https://untappd.com/rss/brewery/462497"

	untappdImageFolder = "../static/images/untappd"

	contentBaseFolder = "../content/beers/"
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

var (
	beerIDRegex   = regexp.MustCompile(`\/b\/[a-z\-0-9]+\/(\d+)`)
	beerNameRegex = regexp.MustCompile(`\/b\/([a-z\-0-9]+)\/[\d]+`)
)

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
	beerID, beerName, beerLink := extractBeerIDFromCheckin(doc)

	contentFolder := filepath.Join(contentBaseFolder, beerName+"-"+beerID)
	if err := os.MkdirAll(contentFolder, 0776); err != nil {
		log.Fatalf("Failed to create content folder %s: %s", contentFolder, err)
	}
	scrapeBasicBeerInfos(beerID, beerLink, contentFolder)

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
		outPath := filepath.Join(contentFolder, imageName)
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

func extractBeerIDFromCheckin(doc *goquery.Document) (beerID, beerName, beerLink string) {
	doc.Find(".checkin-info.pad-it a.label").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("href")
		if exists {
			beerLink = val
			if m := beerIDRegex.FindStringSubmatch(val); len(m) > 1 {
				beerID = m[1]
			} else {
				log.Printf("Failed to extract beerID from link %s", val)
			}

			if m := beerNameRegex.FindStringSubmatch(val); len(m) > 1 {
				beerName = m[1]
			} else {
				log.Printf("Failed to extract beer name from link %s", val)
			}
		}
	})
	return
}

var indexTmpl = template.Must(template.New("index").Parse(`+++
title = "{{.BeerName}}"
description = "{{.Description}}"
abv = "{{.ABV}}"
ibu = "{{.IBU}}"
untappdId = "{{.UntappdID}}"
author = "StackOverflow Brewery"
date = "{{.Date}}"
tags = ["Bier"]
categories = ["Bier"]
comments = false
removeBlur = false
draft = false
+++`))

var (
	detailRegex = regexp.MustCompile(`([\d\.]+).*`)
)

func scrapeBasicBeerInfos(untappdID, beerLink, contentFolder string) {
	indexPath := filepath.Join(contentFolder, "index.md")
	if _, err := os.Stat(indexPath); os.IsExist(err) {
		log.Printf("Index file %s already exists", indexPath)
		return
	}

	resp, err := http.Get("https://untappd.com" + beerLink)
	if err != nil {
		log.Fatalf("Failed to retrieve beer page: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse beer page")
	}
	beerName := doc.Find("h1").First().Text()
	if beerName == "Checking your browser before accessing untappd.com." {
		log.Printf("Ran into browser check, need to try again")
		return
	}
	beerDetailsABVText := doc.Find("div.details p.abv").First().Text()
	beerDetailsIBUText := doc.Find("div.details p.ibu").First().Text()
	beerDetailsDescriptionText := doc.Find("div.desc div.beer-descrption-read-less").First().Text()
	beerDescription := ""
	if len(beerDetailsDescriptionText) > 10 {
		beerDescription = beerDetailsDescriptionText[:len(beerDetailsDescriptionText)-10]
	}
	beerDescription = strings.ReplaceAll(beerDescription, "\n", "")

	indexFile, err := os.Create(indexPath)
	if err != nil {
		log.Fatalf("Failed to create index file %s: %s", indexPath, err)
	}
	defer indexFile.Close()
	tmplData := map[string]string{
		"BeerName":    beerName,
		"Description": beerDescription,
		"UntappdID":   untappdID,
		"ABV":         extractNumber(beerDetailsABVText),
		"IBU":         extractNumber(beerDetailsIBUText),
		"Date":        time.Now().Format(time.RFC3339),
	}

	if err := indexTmpl.Execute(indexFile, tmplData); err != nil {
		log.Fatalf("Failed to render %s: %s", indexPath, err)
	}
}

func extractNumber(in string) string {
	if m := detailRegex.FindStringSubmatch(in); len(m) > 1 {
		return m[1]
	}
	return ""
}
