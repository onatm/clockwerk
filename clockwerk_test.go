package clockwerk

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type DummyJob1 struct{}

func (d DummyJob1) Run(ctx context.Context) {
	fmt.Println("HEY")
}

type DummyJob2 struct{}

func (d DummyJob2) Run(ctx context.Context) {
	time.Sleep(4 * time.Second)
	fmt.Println("HOO")
}

func TestDo(*testing.T) {
	var job1 DummyJob1
	var job2 DummyJob2

	c := New(context.Background())
	c.Every(1 * time.Second).Do(job1)
	c.Every(1 * time.Second).Do(job2)
	defer c.Stop()
	c.Start()

	time.Sleep(5 * time.Second)
}

func TestDo_ManyJobs(t *testing.T) {
	var job1 DummyJob1

	ctx, cancel := context.WithCancel(context.Background())
	c := New(ctx)
	c.Every(5 * time.Second).Do(job1)
	c.Every(10 * time.Second).Do(job1)
	c.Every(1 * time.Minute).Do(job1)
	c.Sharp(5 * time.Second).Do(job1)
	c.Sharp(10 * time.Second).Do(job1)
	c.Sharp(1 * time.Minute).Do(job1)

	c.Start()
	defer c.Stop()

	time.Sleep(10 * time.Second)
	cancel()

	time.Sleep(1 * time.Second)
}
