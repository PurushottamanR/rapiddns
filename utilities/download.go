package utilities

import (
	"errors"
	"io"
	"net/http"
)

func fetchRawData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error fetching raw data from RapidDNS")
	}
	
	defer resp.Body.Close()
	
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading response")
	}
	return string(rawData), nil
}

