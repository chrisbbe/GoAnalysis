:: Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
:: Use of this source code is governed by a BSD-style license that can
:: be found in the LICENSE file.

echo ** Building linux/amd64 **
set GOARCH=amd64
set GOOS=linux
go build -o build/amd64/GoAnalyzerLinux  GoAnalyzer.go

echo ** Building darwin/amd64 **
set GOARCH=amd64
set GOOS=darwin
go build -o build/amd64/GoAnalyzerMac GoAnalyzer.go

echo ** Building windows/amd64 **
set GOARCH=amd64
set GOOS=windows
go build -o build/amd64/GoAnalyzerWindows.exe GoAnalyzer.go
