# Tachograph Go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/way-platform/tachograph-go)](https://pkg.go.dev/github.com/way-platform/tachograph-go)
[![GoReportCard](https://goreportcard.com/badge/github.com/way-platform/tachograph-go)](https://goreportcard.com/report/github.com/way-platform/tachograph-go)
[![CI](https://github.com/way-platform/tachograph-go/actions/workflows/release.yaml/badge.svg)](https://github.com/way-platform/tachograph-go/actions/workflows/release.yaml)

A Go SDK and CLI tool for working with Tachograph data (.DDD files).

## Specification

This SDK implements parsing of downloaded tachograph data, according to [the requirements for the construction, testing, installation, operation and repair of tachographs and their components](https://eur-lex.europa.eu/eli/reg_impl/2016/799/oj/eng).

## Features (roadmap)

> [!CAUTION]
> This SDK is under active development and not yet ready for widespread use.

- Simple interface:

  - `tachograph.UnmarshalFile` to parse a Tachograph file
  - `tachograph.MarshalFile` to serialize a Tachograph file

- Easy to use CLI tool

  - `tachograph parse [...file]`

- Support for generation 1 and 2 (including v2)

- Protobuf-based data model with high usability and full fidelity

- 100% binary marshal/unmarshal round-trip accuracy

- Anonymization of .DDD files (for test data)

- Optional signature validation

## Alternatives

This SDK draws inspiration from other tachograph SDKs, including:

- [traconiq/tachoparser](https://github.com/traconiq/tachoparser)
- [jugglingcats/tachograph-reader](https://github.com/jugglingcats/tachograph-reader)
