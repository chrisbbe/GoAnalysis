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

## License
The MIT License (MIT)
 
Copyright (c) 2015-2016 Christian Bergum Bergersen
 
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
