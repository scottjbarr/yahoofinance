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
	body := `"GE","General Electric Company Common",26.47,26.10,26.02,-0.45,"-1.70%",26.00,26.24,"7/8/2015","10:06am"`

	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 1)

	expected[0] = Quote{
		Symbol:        "GE",
		Name:          "General Electric Company Common",
		PreviousClose: 26.47,
		Open:          26.10,
		Price:         26.02,
		Change:        -0.45,
		ChangePct:     -1.70,
		DayLow:        26.0,
		DayHigh:       26.24,
		LastTradeDate: "7/8/2015",
		LastTradeTime: "10:06am",
	}

	expect(t, 1, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}

func TestGetDataMultipleSymbols(t *testing.T) {
	quoteStrings := make([]string, 2)
	quoteStrings[0] = `"GE","General Electric Company Common",26.47,26.10,26.02,-0.45,"-1.70%",26.00,26.24,"7/8/2015","10:06am"`
	quoteStrings[1] = `"HHGTTG","Towel Company",125.69,124.48,123.49,-2.20,"-1.75%",123.22,124.64,"7/8/2015","10:06am"`

	body := strings.Join(quoteStrings, "\n")

	server, client := httpMock(200, body)
	defer server.Close()

	symbols := []string{"GE", "HHGTTG"}
	quotes := client.GetQuotes(symbols)

	expected := make([]Quote, 2)

	expected[0] = Quote{
		Symbol:        "GE",
		Name:          "General Electric Company Common",
		PreviousClose: 26.47,
		Open:          26.10,
		Price:         26.02,
		Change:        -0.45,
		ChangePct:     -1.70,
		DayLow:        26.0,
		DayHigh:       26.24,
		LastTradeDate: "7/8/2015",
		LastTradeTime: "10:06am",
	}

	expected[1] = Quote{
		Symbol:        "HHGTTG",
		Name:          "Towel Company",
		PreviousClose: 125.69,
		Open:          124.48,
		Price:         123.49,
		Change:        -2.20,
		ChangePct:     -1.75,
		DayLow:        123.22,
		DayHigh:       124.64,
		LastTradeDate: "7/8/2015",
		LastTradeTime: "10:06am",
	}

	expect(t, 2, len(quotes))
	expect(t, reflect.DeepEqual(expected, quotes), true)
}
