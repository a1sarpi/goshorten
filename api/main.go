package main

import (
	"fmt"
	"github.com/a1sarpi/goshorten/api/data"
	"time"
)

func main() {
	shortener := data.NewShortURL("https://sho.rt", "/r/")

	longURL := "https://example.com/very-long-url-to-be-shortened"
	shortURL, err := shortener.Shorten(longURL)
	if err != nil {
		fmt.Println("Error while shortening URL:", err)
		return
	}

	fmt.Printf("Shortened URL: %s → %s\n", longURL, shortURL)

	err = shortener.SetTimeToLive(shortURL, 15*time.Second)
	if err != nil {
		fmt.Println("Error while setting Time To Live:", err)
	} else {
		fmt.Println("Time To Live: 15 seconds")
	}

	ttl, err := shortener.GetTimeToLive(shortURL)
	if err != nil {
		fmt.Println("Error while getting TTL:", err)
	} else {
		fmt.Printf("Time To Live: %.0f seconds\n", ttl.Seconds())
	}

	for i := 0; i < 3; i++ {
		originalURL, err := shortener.GetOriginalURL(shortURL)
		if err != nil {
			fmt.Println("Error while getting original URL:", err)
			break
		}
		fmt.Printf("Short to Long URL #%d: %s → %s\n", i+1, shortURL, originalURL)

		shortener.IncrementClickCount(shortURL)

		time.Sleep(1 * time.Second)
	}

	clicks, err := shortener.GetClickCount(shortURL)
	if err != nil {
		fmt.Println("Error getting click stats:", err)
	} else {
		fmt.Printf("Total clicks: %d\n", clicks)
	}

	fmt.Println("Waiting expiring (15 seconds)...")
	time.Sleep(16 * time.Second)

	_, err = shortener.GetOriginalURL(shortURL)
	if err != nil {
		fmt.Println("After URL is expired:", err)
	}

	expiredCount := shortener.PurgeExpired()
	fmt.Printf("Purged expired links: %d\n", expiredCount)

	_, err = shortener.GetOriginalURL(shortURL)
	if err != nil {
		fmt.Println("After purge:", err)
	}
}
