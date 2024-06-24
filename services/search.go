package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const searchTimeout = 1500 * time.Millisecond

type SearchResponse struct {
	IDs []int `json:"ids"`
}

type SearchRequest struct{}

func Search(query string) SearchResponse {
	searchURL := os.Getenv("SEARCH_URL")
	searchResponse := SearchResponse{
		IDs: []int{},
	}

	if searchURL == "" {
		return searchResponse
	}
	params := url.Values{}
	params.Add("query", query)

	fullURL := fmt.Sprintf("%s?%s", searchURL, params.Encode())

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: searchTimeout,
	}

	// Create a new request
	req, err := http.NewRequest("GET", fullURL, http.NoBody)
	if err != nil {
		slog.Info("error creating request:", "error", err.Error())
		return searchResponse
	}

	// Use context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), searchTimeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		slog.Info("error making search:", "error", err.Error())
		return searchResponse
	}

	if resp != nil && resp.Body != nil {
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				slog.Error("error closing response body", "error", err)
			}
		}()
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		slog.Info("non-200 status code received:", "status", strconv.Itoa(resp.StatusCode))
		return searchResponse
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Info("error reading body:", "error", err.Error())

		return searchResponse
	}

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		slog.Info("error parsing json:", "error", err.Error())

		return searchResponse
	}

	return searchResponse
}
