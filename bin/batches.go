package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/dereulenspiegel/go-brewchild"
)

var (
	brewfatherUserID = os.Getenv("BREWFATHER_USERID")
	brewfatherAPIKey = os.Getenv("BREWFATHER_APIKEY")
)

var (
	dataBaseDir = "../data/batches"
)

type batch struct {
	Name       string   `json:"name"`
	ABV        float64  `json:"abv"`
	IBU        int      `json:"ibu"`
	Color      float64  `json:"color"`
	BrewDate   string   `json:"brewDate"`
	Hops       []string `json:"hops"`
	OG         float64  `json:"og"`
	BuGuRation float64  `json:"buGuRatio`
}

func main() {
	bfClient, err := brewchild.New(brewfatherUserID, brewfatherAPIKey)
	if err != nil {
		log.Fatalf("Failed to create brewfather client: %w", err)
	}

	for _, state := range []string{"Completed", "Planning"} {
		exportBatches(bfClient, state)
	}
}

func exportBatches(bfClient *brewchild.Client, state string) {
	outFilePath := filepath.Join(dataBaseDir, state+".json")
	batches, err := bfClient.Batches(brewchild.Status(state), brewchild.Complete(true), brewchild.Limit(10))
	if err != nil {
		log.Fatalf("Failed to retrieve batches from brewfather: %s", err)
	}
	outFile, err := os.Create(outFilePath)
	if err != nil {
		log.Fatalf("Failed to create output file")
	}

	b := make([]batch, len(batches))

	for i, bt := range batches {
		b[i] = batch{
			Name:  bt.Name,
			ABV:   bt.ABV,
			Color: bt.EstimatedColor,
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
	}
	defer outFile.Close()
	if err := json.NewEncoder(outFile).Encode(b); err != nil {
		log.Fatalf("Failed to encode output data")
	}
}
