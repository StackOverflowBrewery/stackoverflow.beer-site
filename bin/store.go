package main

import (
	"fmt"
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

	if b, exists := s.Beers[br.UntappdID]; exists {
		b.ABV = br.ABV
		b.Description = br.Description
		b.IBU = br.IBU
		b.Name = br.Name
		s.Beers[br.UntappdID] = b
	} else {
		s.Beers[br.UntappdID] = br
	}
	return nil
}

func (s *store) AddBatch(b batch) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if b.UntappdID == "" {
		s.UnlinkedBatches[b.Number] = b
	} else {
		if br, exists := s.Beers[b.UntappdID]; !exists {
			br, err = scrapeBeerDetails(b.UntappdID)
			if err != nil {
				return fmt.Errorf("Failed to scrape untappdID %s for batch %d: %w", b.UntappdID, b.Number, err)
			}
			br.Batches = map[int]batch{
				b.Number: b,
			}
			s.Beers[b.UntappdID] = br
		} else {
			if br.Batches == nil {
				br.Batches = map[int]batch{}
			}
			br.Batches[b.Number] = b
		}
	}
	return nil
}
