// Fetch quotes for equities and currencies from Yahoo Finance.
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

// URL of the remote service
const URL string = "http://download.finance.yahoo.com/d/quotes.csv"

// Minimal Quote data
type Quote struct {
	Symbol        string
	Name          string
	LastTrade     float64
	LastTradeDate string
	LastTradeTime string
}

// Client to the remote service
type Client struct {
	// Base URL
	BaseURL string

	// Requests are transported through this client
	HTTPClient *http.Client
}

// Create a Client with reasonable defaults
func CreateClient() *Client {
	// set the request timeout
	timeout := time.Duration(5 * time.Second)

	client := http.Client{
		Timeout: timeout,
	}

	return &Client{
		BaseURL:    URL,
		HTTPClient: &client,
	}
}

// Return Quotes from the service
func (client *Client) GetQuotes(symbols []string) []Quote {
	// get the body from the HTTP service
	body := client.getData(symbols)

	// get the csv data
	rows := parseCsv(body)

	// create a slice to hold each Quote
	quotes := make([]Quote, len(symbols))

	// create a Quote from each row from the CSV
	for i, row := range rows {
		quotes[i] = buildQuote(row)
	}

	return quotes
}

// Build the URL string.
func (client *Client) buildRequestUrl(symbols []string) string {
	parse_request_url, _ := url.Parse(client.BaseURL)
	parse_request_url.RawQuery = buildParameters(symbols).Encode()

	return parse_request_url.String()
}

// Make the HTTP request to Yahoo.
func (client *Client) getData(symbols []string) string {
	// get the full url
	request_url := client.buildRequestUrl(symbols)

	/// get the HTTP response
	response, err := client.HTTPClient.Get(request_url)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	// read the body
	body_byte, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	return string(body_byte)
}

// Format the symbols suitable for the URL
func formatSymbols(symbols []string) string {
	return strings.Join(symbols, "+")
}

// See Yahoo-data.htm for format details.
func buildParameters(symbols []string) url.Values {
	return url.Values{
		"s": {formatSymbols(symbols)},
		"f": {"snl1d1t1"},
	}
}

// Build a Quote
func buildQuote(data []string) Quote {
	q := Quote{}

	q.Symbol = data[0]
	q.Name = data[1]
	q.LastTrade, _ = strconv.ParseFloat(data[2], 64)
	q.LastTradeDate = data[3]
	q.LastTradeTime = data[4]

	return q
}

// Parse the CSV data from the response.
func parseCsv(body string) [][]string {
	reader := strings.NewReader(body)
	csv := csv.NewReader(reader)

	// fx.Comma = '|'
	// fx.FieldsPerRecord = 16
	csv.LazyQuotes = true

	data, _ := csv.ReadAll()

	return data
}
