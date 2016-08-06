# GoAnalysis

[![Build Status](https://travis-ci.org/chrisbbe/GoAnalysis.svg?branch=master)](https://travis-ci.org/chrisbbe/GoAnalysis) [![GoDoc](https://godoc.org/github.com/chrisbbe/GoAnalysis?status.svg)](https://godoc.org/github.com/chrisbbe/GoAnalysis)

Analyse your Go source code to detect for typical common mistakes in Go and for high values of cyclomatic complecity in your code.

## Install

Requierements:
Go must be installed and $GOPATH have to be set correctly.  

`$ go get github.com/chrisbbe/GoAnalysis/analyzer`


## Execution

`$analyzer -dir="$GOPATH/src/github.com/chrisbbe/GoAnalysis"`

Exchange the example dir with the package you want to analyze.



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
Copyright (c) 2015-2016 The GoAnalysis Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
