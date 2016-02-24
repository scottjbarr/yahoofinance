# Yahoo Finance Quotes

A small Go library to retrieve stock quotes from Yahoo Finance.

There is also a command line program that uses the library

## Build

    go get github.com/scottjbarr/yahoofinance/cmd/yahoofinance

## Usage

Command line program example.

    yahoofinance GE MO AAPL C MSFT USDAUD=X GS

    {AAPL Apple Inc. 124.08 3/26/2015 10:45am}
    {C Citigroup, Inc. Common Stock 50.88 3/26/2015 10:45am}
    {GE General Electric Company Common 24.81 3/26/2015 10:45am}
    {GS Goldman Sachs Group, Inc. (The) 186.375 3/26/2015 10:45am}
    {MO Altria Group, Inc. 50.5 3/26/2015 10:45am}
    {MSFT Microsoft Corporation 41.14 3/26/2015 10:45am}
    {USDAUD=X USD/AUD 1.2762 3/26/2015 2:59pm}

## Testing

    go test

... or alternatively

    goconvery

## References

- [Mocking HTTP Responses in Golang](http://keighl.com/post/mocking-http-responses-in-golang/)

## Licence

The MIT License (MIT)

Copyright (c) 2015 Scott Barr

See [LICENSE.md](LICENSE.md)
