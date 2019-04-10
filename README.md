# Sauron

Sauron is an extensible page parser written in Go.

The purpose of Sauron is to enable the easy implementation of page or website parsers, with first-class support for common platforms or sites such as Reddit and Youtube.

[![](https://img.shields.io/badge/Donate-Flattr-red.svg?style=for-the-badge&link=https://flattr.com/@JoshStrobl&link=https://flattr.com/@JoshStrobl)](https://flattr.com/@JoshStrobl)
[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/JoshStrobl/sauron)
![GitHub](https://img.shields.io/github/license/JoshStrobl/sauron.svg?style=for-the-badge)

## Using

To use Sauron in your application, all you need to do is ensure you are importing `github.com/JoshStrobl/sauron`. Then follow our documentation linked above or look at `tests/` for example code.

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