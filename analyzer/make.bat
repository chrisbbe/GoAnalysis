:: Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
:: Use of this source code is governed by a BSD-style license that can
:: be found in the LICENSE file.

echo ** Building linux/amd64 **
set GOARCH=amd64
set GOOS=linux
go build -v -o build/amd64/analyse  GoAnalyzer.go

echo ** Building windows/amd64 **
set GOARCH=amd64
set GOOS=windows
go build -v -o build/amd64/analyse.exe GoAnalyzer.go
