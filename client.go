// Package yahoofinance provides a way to fetch quotes for equities and
// currencies from Yahoo Finance.
package yahoofinance

import (
	"encoding/csv"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// quoteURL of the remote service
const quotesURL string = "http://download.finance.yahoo.com/d/quotes.csv"

// Quote data
type Quote struct {
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	PreviousClose float64 `json:"previous_close"`
	Open          float64 `json:"open"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePct     float64 `json:"change_pct"`
	DayLow        float64 `json:"day_low"`
	DayHigh       float64 `json:"day_high"`
	LastTradeDate string  `json:"last_trade_date"`
	LastTradeTime string  `json:"last_trade_time"`
	LastTrade     float64 `json:"last_trade"`
}

// Client provides access to the Quotes service.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient returns a new default Client
func NewClient() *Client {
	return &Client{
		BaseURL:    quotesURL,
		HTTPClient: defaultHTTPClient(),
	}
}

// defaultHTTPClient returns http.Client with reasonable defaults
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
}

// GetQuotes return Quotes for the given symbols
func (client *Client) GetQuotes(symbols []string) ([]Quote, error) {
	// get the body from the HTTP service
	body, err := client.get(symbols)

	if err != nil {
		return nil, err
	}

	// get the csv data
	rows := parseCSV(body)

	// create a slice to hold each Quote
	quotes := make([]Quote, len(symbols))

	// create a Quote from each row from the CSV
	for i, row := range rows {
		quotes[i] = buildQuote(row)
	}

	return quotes, nil
}

// buildURL returns a url that can be used to request quotes for the
// given symbols.
func (client *Client) buildURL(symbols []string) string {
	u, _ := url.Parse(client.BaseURL)
	u.RawQuery = buildParameters(symbols).Encode()

	return u.String()
}

// get makes the HTTP request to the service.
func (client *Client) get(symbols []string) (string, error) {
	// get the full url
	url := client.buildURL(symbols)

	// get the HTTP response
	response, err := client.HTTPClient.Get(url)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	// read the body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// formatSymbols returns symbols formatted for the query.
func formatSymbols(symbols []string) string {
	return strings.Join(symbols, "+")
}

// buildParamaters consttucts parameters for the service.
//
// See Yahoo-data.htm for format details.
func buildParameters(symbols []string) url.Values {
	return url.Values{
		"s": {formatSymbols(symbols)},
		"f": {"snpol1c1p2ghd1t1l1"},
	}
}

// parseFloat returns a float64 from the value, handling "%" characters. If
// a value cannot be extract zero is returned.
func parseFloat(value string) float64 {
	if value == "N/A" {
		return 0
	}

	f, err := strconv.ParseFloat(strings.Replace(value, "%", "", 1), 64)

	if err != nil {
		return 0
	}

	return f
}

// buildQuote returns a quote from the array of data
func buildQuote(data []string) Quote {
	q := Quote{}

	q.Symbol = data[0]
	q.Name = data[1]

	q.PreviousClose = parseFloat(data[2])
	q.Open = parseFloat(data[3])
	q.Price = parseFloat(data[4])
	q.Change = parseFloat(data[5])
	q.ChangePct = parseFloat(data[6])
	q.DayLow = parseFloat(data[7])
	q.DayHigh = parseFloat(data[8])

	q.LastTradeDate = data[9]
	q.LastTradeTime = data[10]
	q.LastTrade = parseFloat(data[11])

	return q
}

// parseCSV extracts data from the csv body
func parseCSV(body string) [][]string {
	reader := strings.NewReader(body)
	csv := csv.NewReader(reader)
	csv.LazyQuotes = true

	data, _ := csv.ReadAll()

	return data
}
