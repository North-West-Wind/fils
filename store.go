package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/goware/urlx"
)

var (
	sparseIndex map[uint64]struct{} // Used for deduplication when adding new link
	file        *os.File            // Append-only file
	fileMutex   sync.Mutex
	urls        uint64
	cache       *Cache
)

func LoadStore(storePath string) (err error) {
	// Init variables
	sparseIndex = make(map[uint64]struct{})
	cache = NewCache(1000)

	fileMutex.Lock()
	defer fileMutex.Unlock()
	file, err = os.OpenFile(storePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	urls = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls++
		line := scanner.Text()
		url, err := urlx.NormalizeString(line)
		if err != nil {
			ErrLog.Printf("Failed to normalize %s: %v", line, err)
			continue
		}
		sparseIndex[xxhash.Sum64String(url)] = struct{}{}
	}
	log.Printf("Processed %d lines from store", urls)
	return nil
}

func LookupID(id uint64) (string, error) {
	entry := cache.Get(id)
	if entry != nil {
		return entry.url, nil
	} else {
		// Need to find from file
		fileMutex.Lock()
		defer fileMutex.Unlock()
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			return "", err
		}

		var lineNum uint64
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNum++
			if lineNum == id {
				line := scanner.Text()
				url, err := urlx.NormalizeString(line)
				if err != nil {
					return "", err
				}
				cache.Set(id, url)
				return url, nil
			}
		}

		return "", nil
	}
}

func AddURL(url string) (uint64, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	if _, exists := sparseIndex[xxhash.Sum64String(url)]; exists {
		// URL may already exist. Scan file for it
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			return 0, err
		}

		var lineNum uint64
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNum++
			if scanner.Text() == url {
				return lineNum, nil
			}
		}
	}
	// Seek not needed as we have O_APPEND
	_, err := file.WriteString(url + "\n")
	if err != nil {
		return 0, err
	}
	urls++
	return urls, nil
}
