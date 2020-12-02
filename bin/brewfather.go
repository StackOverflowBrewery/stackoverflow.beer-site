package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/dereulenspiegel/go-brewchild"
)

var (
	brewfatherUserID = os.Getenv("BREWFATHER_USERID")
	brewfatherAPIKey = os.Getenv("BREWFATHER_APIKEY")
)

var (
	dataBaseDir  = "../data/batches"
	batchPostDir = "../content/batches"
)

var (
	untappdIDRegex = regexp.MustCompile(`untappd\(([\d]+)\)`)
)

type malt struct {
	Name     string  `json:"name"`
	Supplier string  `json:"supplier"`
	Amount   float64 `json:"amount"`
}

type misc struct {
	Name   string  `json:"name" yaml:"name"`
	Amount float64 `json:"amount" yaml:"amount"`
	Unit   string  `json:"unit" yaml:"unit"`
	Usage  string  `json:"usage" yaml:"usage"`
}

type hop struct {
	Name   string  `json:"name" yaml:"name"`
	Origin string  `json:"origin" yaml:"origin"`
	Amount float64 `json:"amount" yaml:"amount"`
	Usage  string  `json:"usage" yaml:"usage"`
	Alpha  float64 `json:"alpha" yaml:"alpha"`
}

type batch struct {
	Name        string  `json:"name"`
	ABV         float64 `json:"abv"`
	IBU         int     `json:"ibu"`
	Color       float64 `json:"color"`
	BrewDate    string  `json:"brewDate"`
	OG          float64 `json:"og"`
	BuGuRation  float64 `json:"buGuRatio"`
	UntappdLink string  `json:"untappdLink"`
	UntappdID   string  `json:"untappdId"`
	Number      int     `json:"number"`
	State       string  `json:"state"`
	Malts       []*malt `json:"malts"`
	Hops        []*hop  `json:"hops"`
	Miscs       []*misc `json:"miscs"`
	FG          float64 `json:"fg"`
}

func exportBatches(batches []*brewchild.Batch, state string) []batch {
	outFilePath := filepath.Join(dataBaseDir, state+".json")
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("Failed to create output file")
	}
	defer outFile.Close()

	b := make([]batch, len(batches))

	for i, bt := range batches {
		b[i] = batch{
			Name:   bt.Name,
			ABV:    bt.ABV,
			Color:  bt.EstimatedColor,
			Number: bt.BatchNumber,
			State:  state,
			FG:     bt.MeasuredFG, // TODO convert to plato
		}
		if bt.OG != 0.0 {
			b[i].OG = brewchild.SGToPlato(bt.OG)
		} else {
			b[i].OG = brewchild.SGToPlato(bt.EstimatedOG)
		}
		if b[i].ABV == 0.0 {
			b[i].ABV = bt.MeasuredABV
		}

		if bt.IBU != 0 {
			b[i].IBU = bt.IBU
		} else {
			b[i].IBU = bt.EstimatedIBU
		}

		if bt.BrewDate != nil {
			b[i].BrewDate = bt.BrewDate.String()
		}
		if bt.BuGuRatio != 0.0 {
			b[i].BuGuRation = bt.BuGuRatio
		} else {
			b[i].BuGuRation = bt.EstimatedBuGuRatio
		}

		for _, bh := range bt.Hops {
			b[i].Hops = append(b[i].Hops, &hop{
				Name:   bh.Name,
				Amount: bh.Amount,
				Alpha:  bh.Alpha,
				Origin: bh.Origin,
				Usage:  bh.Usage,
			})
		}

		for _, bm := range bt.Fermentables {
			b[i].Malts = append(b[i].Malts, &malt{
				Name:     bm.Name,
				Supplier: bm.Supplier,
				Amount:   bm.AmountKG,
			})
		}

		for _, bm := range bt.BatchMiscs {
			b[i].Miscs = append(b[i].Miscs, &misc{
				Name:   bm.Name,
				Amount: bm.Amount,
				Unit:   bm.Unit,
				Usage:  bm.Use,
			})
		}

		var untappdID string
		if bt.BatchNotes != "" {
			if m := untappdIDRegex.FindStringSubmatch(bt.BatchNotes); len(m) > 1 {
				untappdID = m[1]
				b[i].UntappdLink = fmt.Sprintf("https://untappd.com/qr/beer/%s", untappdID)
				b[i].UntappdID = m[1]
			}
		}
		ensureBatchPostData(b[i])
	}

	if err := json.NewEncoder(outFile).Encode(b); err != nil {
		log.Fatalf("Failed to encode output data")
	}
	return b
}

var beerContentBase = "../content/beers/"
var batchDataFileName = "batches.json"

func addDataToBeer(b batch, untappdID string) {
	if untappdID == "" {
		log.Printf("UntappdID is empty for batch %d\n", b.Number)
		return
	}

	beerPath, exists := beerContentPath(untappdID)
	if !exists {
		log.Printf("Found no matches for untappdID %s", untappdID)
		return
	}

	batchListPath := filepath.Join(beerPath, batchDataFileName)

	batches := []string{}
	if _, err := os.Stat(batchListPath); os.IsNotExist(err) {
		batches = append(batches, fmt.Sprintf("%d", b.Number))
		outFile, err := os.Create(batchListPath)
		if err != nil {
			log.Fatalf("Failed to create batch list for beer: %s", err)
		}
		defer outFile.Close()
		if err := json.NewEncoder(outFile).Encode(&batches); err != nil {
			log.Fatalf("Failed to marshal batch list: %s", err)
		}
	} else {
		outFile, err := os.OpenFile(batchListPath, os.O_RDWR, 0776)
		if err != nil {
			log.Fatalf("Failed to open existing batch file %s: %s", batchListPath, err)
		}
		defer outFile.Close()
		if err := json.NewDecoder(outFile).Decode(&batches); err != nil {
			log.Fatalf("Failed to unmarshal existing batch list: %s", err)
		}

		currBatch := fmt.Sprintf("%d", b.Number)
		exists := false
		for _, bn := range batches {
			if bn == currBatch {
				exists = true
			}
		}

		if !exists {
			batches = append(batches, currBatch)
			if err := json.NewEncoder(outFile).Encode(&batches); err != nil {
				log.Fatalf("Failed to marshal batch list: %s", err)
			}
		}
	}
}

func ensureBatchPostData(b batch) {
	batchFolder := filepath.Join(batchPostDir, fmt.Sprintf("%d", b.Number))
	if err := os.MkdirAll(batchFolder, 0776); err != nil {
		log.Fatalf("Unable to create batch folder %s: %s", batchFolder, err)
	}
	dataFilePath := filepath.Join(batchFolder, "data.json")

	if _, err := os.Stat(dataFilePath); os.IsExist(err) {
		os.Remove(dataFilePath)
	}
	dataFile, err := os.Create(dataFilePath)
	if err != nil {
		log.Fatalf("Failed to create batch data file %s: %s", dataFilePath, err)
	}
	defer dataFile.Close()
	if err := json.NewEncoder(dataFile).Encode(b); err != nil {
		log.Fatalf("Failed to marshal batch data file %s: %s", dataFilePath, err)
	}
}
