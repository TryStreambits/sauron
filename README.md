# Sauron

Sauron is an extensible page parser written in Go.

The purpose of Sauron is to enable the easy implementation of page or website parsers, with first-class support for common platforms or sites such as Reddit and Youtube.

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/TryStreambits/sauron)
[![goreportcard](https://img.shields.io/badge/goreportcard-A+-green.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/TryStreambits/sauron)
![GitHub](https://img.shields.io/github/license/TryStreambits/sauron.svg?style=for-the-badge)

## Using

To use Sauron in your application, all you need to do is ensure you are importing `github.com/TryStreambits/sauron`. Then follow the documentation linked above or look at `tests/` for example code.

## Building

To compile, first ensure you have turned on Go Module support if you are working inside your `GOPATH`:

``` bash
export GO111MODULE=on
```

Next, all you have to do is run the following command to compile:

``` bash
go build
```

## License

Sauron is licensed under the Apache-2.0 license.