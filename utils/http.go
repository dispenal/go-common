package common_utils

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Fetch(url string, headers ...map[string]string) (*goquery.Document, error) {
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	for _, header := range headers {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}

	// Create a new http client
	client := &http.Client{}

	// Send the request via a client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Defer the closing of the body
	defer res.Body.Close()

	data, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
