package download

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent 64 possibilities (62 actually used)
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func FetchJson(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	return body, nil
}

func GenerateSeed(length int) string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	b := make([]byte, length)
	for i := range length {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}
