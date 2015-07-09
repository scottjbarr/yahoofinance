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
		"f": {"snpol1c1p2ghd1t1l1"},
	}
}

// Parse a float64 from the value, or return zero
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

// Build a Quote
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
