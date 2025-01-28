package marketdata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AlphaVantage struct {
	apiKey     string
	httpClient *http.Client
}

type QuoteResponse struct {
	GlobalQuote struct {
		Symbol        string `json:"01. symbol"`
		Price         string `json:"05. price"`
		Volume        string `json:"06. volume"`
		LatestDay     string `json:"07. latest trading day"`
		PreviousClose string `json:"08. previous close"`
	} `json:"Global Quote"`
}

func NewAlphaVantage(apiKey string) *AlphaVantage {
	return &AlphaVantage{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (av *AlphaVantage) GetQuote(symbol string) (*QuoteResponse, error) {
	url := fmt.Sprintf(
		"https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s",
		symbol,
		av.apiKey,
	)

	resp, err := av.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quoteResp QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResp); err != nil {
		return nil, err
	}

	return &quoteResp, nil
}
