# Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style license that can
# be found in the LICENSE file.

PLATFORMS := linux/amd64 darwin/amd64 windows/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) go build -o 'GoAnalyzer-$(os)-$(arch)' GoAnalyzer.go

#.PHONY release $(PLATFORMS)

clean:
	if [ -f GoAnalyzer-linux-amd64 ] ; then rm GoAnalyzer-linux-amd64 ; fi
	if [ -f GoAnalyzer-darwin-amd64 ] ; then rm GoAnalyzer-darwin-amd64 ; fi
	if [ -f GoAnalyzer-windows-amd64 ] ; then rm GoAnalyzer-windows-amd64 ; fi
