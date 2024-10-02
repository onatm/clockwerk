# clockwerk

[![Coverage Status](https://coveralls.io/repos/github/onatm/clockwerk/badge.svg?branch=main)](https://coveralls.io/github/onatm/clockwerk?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/onatm/clockwerk)](https://goreportcard.com/report/github.com/onatm/clockwerk)
[![GoDoc](http://godoc.org/github.com/onatm/clockwerk?status.png)](http://godoc.org/github.com/onatm/clockwerk) 

Job Scheduling Library

clockwerk allows you to schedule periodic jobs using a simple, fluent syntax.

## Installing

Using clockwerk is easy. First, use `go get` to install the latest version of the library.

``` sh
go get -u github.com/onatm/clockwerk@latest
```

## Usage

Include clockwerk in your application:

```go
import "github.com/onatm/clockwerk"
```

## Example

``` go
package main

import (
  "fmt"
  "time"
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
