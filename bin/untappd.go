package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
	contentFolder := filepath.Join(contentBaseFolder, beerName+"-"+beerID)
	if err := os.MkdirAll(contentFolder, 0776); err != nil {
		return "", fmt.Errorf("Failed to create beer content folder %s: %w", contentFolder, err)
	}
	return contentFolder, nil
}
