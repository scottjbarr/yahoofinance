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
	body := `"GE","General Electric Company Common",26.03,"6/30/2015","10:43am",-0.28,"-1.06%"`
	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 1)

	expected[0] = Quote{
		Symbol:        "GE",
		Name:          "General Electric Company Common",
		LastTrade:     26.03,
		LastTradeDate: "6/30/2015",
		LastTradeTime: "10:43am",
		Change:        -0.28,
		ChangePct:     -1.06,
	}

	expect(t, 1, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}

func TestGetDataMultipleSymbols(t *testing.T) {
	quoteStrings := make([]string, 2)
	quoteStrings[0] = `"GE","General Electric Company Common",26.025,"6/30/2015","10:44am",-0.285,"-1.083%"`
	quoteStrings[1] = `"HHGTTG","Towel Company",124.42,"6/30/2015","10:45am",-1.58,"-1.25%"`
	body := strings.Join(quoteStrings, "\n")

	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE", "HHGTTG"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 2)

	expected[0] = Quote{
		Symbol:        "GE",
		Name:          "General Electric Company Common",
		LastTrade:     26.025,
		LastTradeDate: "6/30/2015",
		LastTradeTime: "10:44am",
		Change:        -0.285,
		ChangePct:     -1.083,
	}

	expected[1] = Quote{
		Symbol:        "HHGTTG",
		Name:          "Towel Company",
		LastTrade:     124.42,
		LastTradeDate: "6/30/2015",
		LastTradeTime: "10:45am",
		Change:        -1.58,
		ChangePct:     -1.25,
	}

	expect(t, 2, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}
