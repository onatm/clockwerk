package clockwerk

import (
	"fmt"
	"testing"
	"time"
)

type DummyJob1 struct{}

func (d DummyJob1) Run() {
	fmt.Println("HEY")
}

type DummyJob2 struct{}

func (d DummyJob2) Run() {
	time.Sleep(4 * time.Second)
	fmt.Println("HOO")
}

func TestDo(*testing.T) {
	var job1 DummyJob1
	var job2 DummyJob2

	c := New()
	c.Every(1 * time.Second).Do(job1)
	c.Every(1 * time.Second).Do(job2)
	defer c.Stop()
	c.Start()

	time.Sleep(5 * time.Second)
}
