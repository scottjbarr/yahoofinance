package yahoofinance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// Mock out HTTP requests.
//
// Pinched from http://keighl.com/post/mocking-http-responses-in-golang/
// Thanks, Kyle Truscott (@keighl)!
func httpMock(code int, body string) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/octet-stream")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &Client{BaseURL: server.URL, HTTPClient: httpClient}

	return server, client
}

// Test helper. Thanks again, @keighl
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestGetData(t *testing.T) {
	body := `"GE","General Electric Company Common",24.91,"3/25/2015","4:00pm"`
	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 1)

	expected[0] = Quote{
		Symbol:        "GE",
		Name:          "General Electric Company Common",
		LastTrade:     24.91,
		LastTradeDate: "3/25/2015",
		LastTradeTime: "4:00pm",
	}

	expect(t, 1, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}

func TestGetDataMultipleSymbols(t *testing.T) {
	quoteStrings := make([]string, 2)
	quoteStrings[0] = `"ACME","ACME Company",20.15,"3/25/2015","4:00pm"`
	quoteStrings[1] = `"HHGTTG","Towel Company",42.01,"3/26/2015","10:12am"`
	body := strings.Join(quoteStrings, "\n")

	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE", "FOO"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 2)

	expected[0] = Quote{
		Symbol:        "ACME",
		Name:          "ACME Company",
		LastTrade:     20.15,
		LastTradeDate: "3/25/2015",
		LastTradeTime: "4:00pm",
	}

	expected[1] = Quote{
		Symbol:        "HHGTTG",
		Name:          "Towel Company",
		LastTrade:     42.01,
		LastTradeDate: "3/26/2015",
		LastTradeTime: "10:12am",
	}

	expect(t, 2, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}
