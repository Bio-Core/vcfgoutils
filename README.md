vcfgoutils - a Go library for handling Variant Call Format (VCF) files
===

[![Go Report Card](https://goreportcard.com/badge/github.com/bio-core/vcfgoutils)](https://goreportcard.com/report/github.com/bio-core/vcfgoutils)
[![Build Status](https://travis-ci.org/Bio-Core/vcfgoutils.svg?branch=master)](https://travis-ci.org/Bio-Core/vcfgoutils)

`vcfgoutils` is a `golang` package used to import, process and convert VCF
files into a different format for downstream analysis.

### Prerequisites
The following go packages are required to run the vcfgoutils suite of tools.
* github.com/brentp/vcfgo
* github.com/nats-io/nats
* gopkg.in/mgo.v2

In addition, tools are available to upload converted VCF data to a MongoDB
server.

### Installing
This package can be downloaded from gitlab.com using:

```
go get gitlab.com/uhn/vcfgoutils
```

You can subsequently test, build and install the package using:

```
go test
go build
go install
```

### License
MIT License

