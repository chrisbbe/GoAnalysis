# GoAnalysis

[![Build Status](https://travis-ci.org/chrisbbe/GoAnalysis.svg?branch=master)](https://travis-ci.org/chrisbbe/GoAnalysis)

Analyse your Go source code to detect for typical common mistakes in Go and for high values of cyclomatic complecity in your code.

## Tests

GoAnalysis is developed using the Test-driven development (TDD) process where unit tests are extensively used to guarantee for the functionality in each package and the hole analysis as a unit.

Run the following command in the root folder to execute all tests in all packages.

`$ go test ./...`

or just the following command at root level in each package you want to test.

`$ go test`

## Code Style

### Indentation

**Go**
```
Tab size: 2
Indent: 2
Continuous indent: 2
```
