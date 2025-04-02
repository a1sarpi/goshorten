package data

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

const (
	allowedLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	urlLength = 10
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ShortURL struct {
	mu          sync.RWMutex
	host        string
	path        string
	shortToLong map[string]string
	longToShort map[string]string
	clicksCount map[string]int
	expires     map[string]time.Time
}

func NewShortURL(host, path string) *ShortURL {
	return &ShortURL{
		host:        host,
		path:        path,
		shortToLong: make(map[string]string),
		longToShort: make(map[string]string),
		clicksCount: make(map[string]int),
		expires:     make(map[string]time.Time),
	}
}

func (sh *ShortURL) Shorten(longURL string) (string, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	if shortURL, exists := sh.longToShort[longURL]; exists {
		return shortURL, nil
	}

	shortCode := generateShortCode()
	for exists := true; exists; _, exists = sh.shortToLong[sh.buildShortURL(shortCode)] {
		shortCode = generateShortCode()
	}
	shortURL := sh.buildShortURL(shortCode)

	sh.shortToLong[shortURL] = longURL
	sh.longToShort[longURL] = shortURL
	sh.clicksCount[shortURL] = 0

	return shortURL, nil
}

func (sh *ShortURL) GetOriginalURL(shortURL string) (string, error) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	longURL, exists := sh.shortToLong[shortURL]
	if !exists {
		return "", ErrURLNotFound
	}

	if expiry, exists := sh.expires[shortURL]; exists && time.Now().After(expiry) {
		return "", ErrURLNotFound
	}

	return longURL, nil
}

func (sh *ShortURL) Delete(shortURL string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	longURL, exists := sh.shortToLong[shortURL]
	if !exists {
		return ErrURLNotFound
	}

	delete(sh.longToShort, longURL)
	delete(sh.shortToLong, shortURL)
	delete(sh.clicksCount, shortURL)
	delete(sh.expires, shortURL)

	return nil
}

func (sh *ShortURL) GetClickCount(shortURL string) (int, error) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	count, exists := sh.clicksCount[shortURL]
	if !exists {
		return -1, ErrURLNotFound
	}

	return count, nil
}

func (sh *ShortURL) IncrementClickCount(shortURL string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	sh.clicksCount[shortURL]++
}

func (sh *ShortURL) SetTimeToLive(shortURL string, duration time.Duration) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	if _, exists := sh.shortToLong[shortURL]; !exists {
		return ErrURLNotFound
	}

	sh.expires[shortURL] = time.Now().Add(duration)
	return nil
}

func (sh *ShortURL) GetTimeToLive(shortURL string) (time.Duration, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	expiry, exists := sh.expires[shortURL]
	if !exists {
		return 0, ErrURLNotFound
	}

	return time.Until(expiry), nil
}

func (sh *ShortURL) PurgeExpired() int {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	count := 0
	now := time.Now()

	for shortURL, expiry := range sh.expires {
		if now.After(expiry) {
			longURL := sh.shortToLong[shortURL]
			delete(sh.longToShort, longURL)
			delete(sh.shortToLong, shortURL)
			delete(sh.clicksCount, shortURL)
			delete(sh.expires, shortURL)
			count++
		}
	}

	return count
}

func generateShortCode() string {
	b := make([]byte, urlLength)
	for i := range b {
		b[i] = allowedLetters[rnd.Intn(len(allowedLetters))]
	}

	return string(b)
}

func (sh *ShortURL) buildShortURL(shortCode string) string {
	var builder strings.Builder
	builder.Grow(urlLength + len(sh.host) + len(sh.path))

	builder.WriteString(sh.host)
	builder.WriteString(sh.path)
	builder.WriteString(shortCode)

	return builder.String()
}
