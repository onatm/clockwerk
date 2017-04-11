# clockwerk
Job Scheduling Library [![Build Status](https://travis-ci.org/onatm/clockwerk.svg?branch=master)](https://travis-ci.org/onatm/clockwerk) [![GoDoc](http://godoc.org/github.com/onatm/clockwerk?status.png)](http://godoc.org/github.com/onatm/clockwerk) 

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
  c.EverySeconds(30).Do(job)
  c.Start()
}
```