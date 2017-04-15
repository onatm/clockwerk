# clockwerk

[![Build Status](https://travis-ci.org/onatm/clockwerk.svg?branch=master)](https://travis-ci.org/onatm/clockwerk) &nbsp; [![Coverage Status](https://coveralls.io/repos/github/onatm/clockwerk/badge.svg?branch=master)](https://coveralls.io/github/onatm/clockwerk?branch=master) &nbsp; [![Go Report Card](https://goreportcard.com/badge/github.com/onatm/clockwerk)](https://goreportcard.com/report/github.com/onatm/clockwerk) &nbsp; [![GoDoc](http://godoc.org/github.com/onatm/clockwerk?status.png)](http://godoc.org/github.com/onatm/clockwerk) 

Job Scheduling Library

clockwerk allows you to schedule periodic jobs using a simple, fluent syntax.

## Usage

``` sh
go get github.com/onatm/clockwerk
```

``` go
package main

import (
  "fmt"
  "github.com/onatm/clockwerk"
)

type DummyJob struct{}

func (d DummyJob) Run() {
  fmt.Println("Every 30 seconds")
}

func main() {
  var job DummyJob
  c := clockwerk.New()
  c.Every(30 * time.Second).Do(job)
  c.Start()
}
```
