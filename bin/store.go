package main

import (
	"fmt"
	"log"
	"sync"
)

type beer struct {
	Name        string
	Description string
	UntappdID   string
	ABV         string
	IBU         string
	Batches     map[int]batch
}

func (b *beer) LatestBatch() *batch {
	currNum := 0
	var currBatch *batch
	for num, bt := range b.Batches {
		if num > currNum {
			currNum = num
			currBatch = &bt
		}
	}
	return currBatch
}

type store struct {
	Beers           map[string]*beer
	UnlinkedBatches map[int]batch

	lock *sync.Mutex
}

var (
	str = &store{
		Beers:           map[string]*beer{},
		UnlinkedBatches: map[int]batch{},
		lock:            &sync.Mutex{},
	}
)

func (s *store) AddBeer(br *beer) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if br.UntappdID == "" {
		return fmt.Errorf("Beer %s has no valid untappdID", br.Name)
	}

	if b, exists := s.Beers[br.UntappdID]; exists {
		log.Printf("Updating beer %s (%s)", b.Name, b.UntappdID)
		b.ABV = br.ABV
		b.Description = br.Description
		b.IBU = br.IBU
		b.Name = br.Name
		s.Beers[br.UntappdID] = b
	} else {
		log.Printf("Adding new beer %s (%s)", br.Name, br.UntappdID)
		s.Beers[br.UntappdID] = br
	}
	return nil
}

func (s *store) AddBatch(b batch) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if b.UntappdID == "" {
		log.Printf("Batch %d has no untappd info associated", b.Number)
		s.UnlinkedBatches[b.Number] = b
	} else {
		if br, exists := s.Beers[b.UntappdID]; !exists {
			log.Printf("Batch %d belongs non scraped beer %s, scraping...", b.Number, b.UntappdID)
			br, err = scrapeBeerDetails(b.UntappdID)
			if err != nil {
				return fmt.Errorf("Failed to scrape untappdID %s for batch %d: %w", b.UntappdID, b.Number, err)
			}
			br.Batches = map[int]batch{
				b.Number: b,
			}
			s.Beers[b.UntappdID] = br
		} else {
			log.Printf("Batch %d belong to already scraped beer %s (%s)", b.Number, br.Name, b.UntappdID)
			if br.Batches == nil {
				br.Batches = map[int]batch{}
			}
			br.Batches[b.Number] = b
		}
	}
	return nil
}
