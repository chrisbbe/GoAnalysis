// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import (
	"bytes"
	"fmt"
	"log"
)

// Trivial type to hold address.
type Address struct {
	Street, City string
}

type Person struct {
	FirstName, LastName string
}

type T1 string
type T2 string
type T3 string

var logger *log.Logger

func main() {
	myAddress := Address{
		Street: "Trimveien 6",
		City:   "Oslo",
	}
	log.Println(myAddress)

	me := Person{
		FirstName: "Christian",
		LastName:  "Bergum Bergersen",
	}
	log.Println(me)

	var foo T1
	foo = "foo"
	log.Println(foo)

	var bar T3
	bar = "bar"
	log.Println(bar)
}

// Calls itself.
func (address Address) String() string {
	return fmt.Sprintf("%s", address)
}

// Correct.
func (person Person) String() string {
	return fmt.Sprintf("%s %s", person.FirstName, person.LastName)
}

// Calls itself.
func (t1 T1) String() string {
	return fmt.Sprint(t1)
}

// Calls itself.
func (t2 T2) String() string {
	log.Print("Calling String() for : %v", t2)
	return fmt.Sprintln(t2)
}

// Calls itself.
func (bar T3) String() string {
	var buf bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	logger.Printf("Calling String() for %+v", bar)
	return fmt.Sprintf("Bar")
}
