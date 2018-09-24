Stow 
==========

[![GoDoc](https://godoc.org/github.com/djherbis/stow?status.svg)](https://godoc.org/github.com/djherbis/stow)
[![Release](https://img.shields.io/github/release/djherbis/stow.svg)](https://github.com/djherbis/stow/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.txt)
[![Build Status](https://travis-ci.org/djherbis/stow.svg?branch=master)](https://travis-ci.org/djherbis/stow) 
[![Coverage Status](https://coveralls.io/repos/djherbis/stow/badge.svg?branch=master)](https://coveralls.io/r/djherbis/stow?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/djherbis/stow)](https://goreportcard.com/report/github.com/djherbis/stow)

Usage
------------

This package provides a persistence manager for objects backed by boltdb.

```go
package main

import (
  "encoding/gob"
  "fmt"
  "log"

  "github.com/boltdb/bolt"
  "gopkg.in/djherbis/stow.v2"
)

func main() {
  // Create a boltdb database
  db, err := bolt.Open("my.db", 0600, nil)
  if err != nil {
    log.Fatal(err)
  }

  // Open/Create a Json-encoded Store, Xml and Gob are also built-in
  // We'll we store a greeting and person in a boltdb bucket named "people"
  peopleStore := stow.NewJSONStore(db, []byte("people"))

  peopleStore.Put("hello", Person{Name: "Dustin"})

  peopleStore.ForEach(func(greeting string, person Person) {
    fmt.Println(greeting, person.Name)
  })

  // Open/Create a Gob-encoded Store. The Gob encoding keeps type information,
  // so you can encode/decode interfaces!
  sayerStore := stow.NewStore(db, []byte("greetings"))

  var sayer Sayer = Person{Name: "Dustin"}
  sayerStore.Put("hello", &sayer)

  var retSayer Sayer
  sayerStore.Get("hello", &retSayer)
  retSayer.Say("hello")

  sayerStore.ForEach(func(sayer Sayer) {
    sayer.Say("hey")
  })
}

type Sayer interface {
  Say(something string)
}

type Person struct {
  Name string
}

func (p Person) Say(greeting string) {
  fmt.Printf("%s says %s.\n", p.Name, greeting)
}

func init() {
  gob.Register(&Person{})
}

```

Installation
------------
```sh
go get gopkg.in/djherbis/stow.v2
```
