# gtfs #

[![Build Status](https://github.com/dpearson/gtfs/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/dpearson/gtfs/actions)
[![GoDoc](https://godoc.org/github.com/dpearson/gtfs?status.svg)](https://godoc.org/github.com/dpearson/gtfs)

This package allows [General Transit Feed Specification](https://developers.google.com/transit/gtfs/) files to be read and manipulated from Go.

In addition to basic GTFS support, this package also supports the following [Google Transit Extensions to GTFS](https://developers.google.com/transit/gtfs/reference/gtfs-extensions):

* [Extended GTFS route types](https://developers.google.com/transit/gtfs/reference/extended-route-types)
* [Station vehicle types](https://developers.google.com/transit/gtfs/reference/gtfs-extensions#station-vehicle-types)
* [Station platforms](https://developers.google.com/transit/gtfs/reference/gtfs-extensions#station-platforms)
* [Trip diversions](https://developers.google.com/transit/gtfs/reference/gtfs-extensions#trip-diversions)
* [Translations](https://developers.google.com/transit/gtfs/reference/gtfs-extensions#translations)

Support for additional extensions may be added in the future.

## Installation ##

All code can be downloaded with:

```sh
go get -u github.com/dpearson/gtfs
```

This package has no dependencies outside of the standard library.

## License ##

```
MIT License

Copyright (c) 2018-2023 David Pearson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```