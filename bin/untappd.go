package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gernest/front"
	"gopkg.in/yaml.v3"
)

var frontMatterDelimiter = []byte("---\n")

type frontMatter struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	ABV         float64 `yaml:"abv"`
	IBU         int     `yaml:"ibu"`
	UntappdID   string  `yaml:"untappdId"`

	Author     string    `yaml:"author"`
	Date       time.Time `yaml:"date"`
	Tags       []string  `yaml:"tags"`
	Categories []string  `yaml:"categories"`
	Comments   bool      `yaml:"comments"`
	RemoveBlur bool      `yaml:"removeBlur"`
	Draft      bool      `yaml:"draft"`
}

var (
	detailRegex = regexp.MustCompile(`([\d\.]+).*`)
)

var (
	breweryFeed = "https://untappd.com/rss/brewery/462497"

	untappdImageFolder = "../static/images/untappd"

	contentBaseFolder = "../content/beers/"
)

var (
	beerIDRegex   = regexp.MustCompile(`\/b\/[a-z\-0-9]+\/(\d+)`)
	beerNameRegex = regexp.MustCompile(`\/b\/([a-z\-0-9]+)\/[\d]+`)
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

func ensureBeerContent(contentFolder string, b *beer) (err error) {
	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	var body string
	var indexFile *os.File

	indexPath := filepath.Join(contentFolder, "index.md")

	if indexFile, err = os.Open(indexPath); err == nil {
		log.Printf("Index file %s already exists, updating front matter", indexPath)
		_, body, err = m.Parse(indexFile)
		if err != nil {
			indexFile.Close()
			return fmt.Errorf("Failed to parse front matter of %s: %w", indexPath, err)
		}
		indexFile.Close()
	}

	fmBuf := bytes.NewBuffer(nil)
	fm := &frontMatter{
		Title:       b.Name,
		Description: b.Description,
		UntappdID:   b.UntappdID,
		Author:      "StackOverflow Brewery",
		Tags:        []string{"Bier"},
		Categories:  []string{"Bier"},
		Comments:    false,
		RemoveBlur:  false,
		Draft:       false,
	}

	if bt := b.LatestBatch(); bt != nil {
		fm.Date, err = time.Parse(time.RFC3339, bt.BrewDate)
		if err != nil {
			return fmt.Errorf("Found unparseable date on batch %d: %w", bt.Number, err)
		}
		fm.ABV = bt.ABV
		fm.IBU = bt.IBU
	} else {
		abv, err := strconv.ParseFloat(b.ABV, 64)
		if err != nil {
			return fmt.Errorf("Failed to parse abv string %s: %w", b.ABV, err)
		}
		ibu, err := strconv.Atoi(b.IBU)
		if err != nil {
			return fmt.Errorf("Failed to parse ibu string %s: %w", b.IBU, err)
		}
		fm.ABV = abv
		fm.IBU = ibu
		fm.Date, _ = time.Parse(time.RFC3339, "2020-01-01T01:00:00+01:00")
	}

	if err := yaml.NewEncoder(fmBuf).Encode(fm); err != nil {
		return fmt.Errorf("Failed to marshal frontmatter: %w", err)
	}

	indexFile, err = os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("Failed to create index file %s: %w", indexPath, err)
	}
	defer indexFile.Close()

	if _, err := indexFile.Write(frontMatterDelimiter); err != nil {
		return fmt.Errorf("Failed to write to %s: %w", indexPath, err)
	}
	if _, err := indexFile.Write(fmBuf.Bytes()); err != nil {
		return fmt.Errorf("Failed to write to %s: %w", indexPath, err)
	}
	if _, err := indexFile.Write([]byte("\n")); err != nil {
		return fmt.Errorf("Failed to write to %s: %w", indexPath, err)
	}
	if _, err := indexFile.Write(frontMatterDelimiter); err != nil {
		return fmt.Errorf("Failed to write to %s: %w", indexPath, err)
	}
	if _, err := indexFile.Write([]byte(body)); err != nil {
		return fmt.Errorf("Failed to write to %s: %w", indexPath, err)
	}

	return nil
}

func scrapeBeerDetails(untappdID string) (*beer, error) {
	log.Printf("Scraping details for beer %s", untappdID)
	beerLink := "/b/awesome-beer/" + untappdID
	resp, err := http.Get("https://untappd.com" + beerLink)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve beer page: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse beer page: %w", err)
	}
	if err := checkBrowserCheck(doc); err != nil {
		return nil, err
	}
	beerName := doc.Find("h1").First().Text()

	beerDetailsABVText := doc.Find("div.details p.abv").First().Text()
	beerDetailsIBUText := doc.Find("div.details p.ibu").First().Text()
	beerDetailsDescriptionText := doc.Find("div.desc div.beer-descrption-read-less").First().Text()
	beerDescription := ""
	if len(beerDetailsDescriptionText) > 10 {
		beerDescription = beerDetailsDescriptionText[:len(beerDetailsDescriptionText)-10]
	}
	beerDescription = strings.ReplaceAll(beerDescription, "\n", "")

	return &beer{
		UntappdID:   untappdID,
		Name:        beerName,
		Description: beerDescription,
		ABV:         extractNumber(beerDetailsABVText),
		IBU:         extractNumber(beerDetailsIBUText),
	}, nil
}

func parseCheckin(checkinURL string) (beerID, imageLink string) {
	resp, err := http.Get(checkinURL)
	if err != nil {
		log.Fatalf("Failed to retrieve checkin page: %s", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse checkin page")
	}
	if err := checkBrowserCheck(doc); err != nil {
		log.Printf("Ran into browser check for checkin %s", checkinURL)
		return
	}

	beerID, beerName, _ := extractBeerIDFromCheckin(doc)

	contentFolder, err := ensureBeerContentFolder(beerID, beerName)
	if err != nil {
		log.Fatalf("Failed to create beer content folder: %s", err)
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("property")
		if exists && val == "og:image" {
			contentVal, exists := s.Attr("content")
			if exists {
				imageLink = contentVal
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
	return beerID, imageLink
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

func beerContentPath(beerID string) (string, bool) {
	files, err := filepath.Glob(contentBaseFolder + "/*" + beerID)
	if err != nil {
		return "", false
	}
	if len(files) > 0 {
		return files[0], true
	}
	return "", false
}

func ensureBeerContentFolder(beerID, beerName string) (string, error) {
	if p, exists := beerContentPath(beerID); exists {
		return p, nil
	}
	beerFolderName := cleanBeerName(beerName)
	contentFolder := filepath.Join(contentBaseFolder, beerFolderName+"-"+beerID)
	if err := os.MkdirAll(contentFolder, 0776); err != nil {
		return "", fmt.Errorf("Failed to create beer content folder %s: %w", contentFolder, err)
	}
	return contentFolder, nil
}

func cleanBeerName(in string) string {
	out := strings.ToLower(in)
	out = strings.ReplaceAll(out, " ", "_")
	out = strings.ReplaceAll(out, ".", "")
	return out
}

func extractNumber(in string) string {
	if m := detailRegex.FindStringSubmatch(in); len(m) > 1 {
		return m[1]
	}
	return ""
}

func checkBrowserCheck(doc *goquery.Document) error {
	beerName := doc.Find("h1").First().Text()
	if beerName == "Checking your browser before accessing untappd.com." {
		return errors.New("Browser check")
	}
	return nil
}
