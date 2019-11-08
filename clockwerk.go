package clockwerk

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// Job is an interface for added jobs.
type Job interface {
	Run(ctx context.Context)
}

// Entry consists of a schedule and the job to execute on that schedule.
type Entry struct {
	Period time.Duration
	Sharp  bool
	Next   time.Time
	Prev   time.Time
	Job    Job
}

func newEntry(period time.Duration, sharp bool) *Entry {
	return &Entry{
		Period: period,
		Prev:   time.Unix(0, 0),
		Next:   time.Unix(0, 0),
		Sharp:  sharp,
	}
}

// Do adds a Job to the Entry.
func (e *Entry) Do(job Job) {
	e.Job = job
}

// Clockwerk keeps track of any number of entries, invoking associated Job's Run
// method as specified by the schedule.
type Clockwerk struct {
	ctx     context.Context
	entries []*Entry
	stop    chan struct{}
	running bool
	mu      sync.Mutex
}

// New returns a new Clockwerk job runner.
func New(ctx context.Context) *Clockwerk {
	return &Clockwerk{
		ctx:     ctx,
		entries: nil,
		stop:    make(chan struct{}),
		running: false,
	}
}

// Every schedules a new Entry and return it.
func (c *Clockwerk) Every(period time.Duration) *Entry {
	return c.newEntry(period, false)
}

// Sharp schedules a new Entry on clock and return it.
func (c *Clockwerk) Sharp(unit time.Duration) *Entry {
	return c.newEntry(unit, true)
}

func (c *Clockwerk) newEntry(period time.Duration, sharp bool) *Entry {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return nil
	}

	entry := newEntry(period, sharp)
	if sharp {
		c.schedule(entry)
	}

	c.entries = append(c.entries, entry)

	return entry
}

// Start the Clockwerk in its own go-routine, or no-op if already started.
func (c *Clockwerk) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// Stop the Clockwerk if it is running.
func (c *Clockwerk) Stop() {
	c.mu.Lock()
	c.mu.Unlock()

	if !c.running {
		return
	}
	close(c.stop)
	c.running = false
}

func (c *Clockwerk) schedule(e *Entry) {
	e.Prev = time.Now()

	if e.Sharp {
		e.Next = e.Prev.Add(e.Period).Truncate(e.Period)
	} else {
		e.Next = e.Prev.Add(e.Period)
	}
}

func (c *Clockwerk) run() {
	for {
		// Sort all entries
		sort.Slice(c.entries, func(i, j int) bool {
			return c.entries[i].Next.Before(c.entries[j].Next)
		})

		t := time.After(c.entries[0].Next.Sub(time.Now()))
		select {
		case <-c.ctx.Done():
			return
		case <-c.stop:
			return
		case <-t:
			c.runPending()
			continue
		}
	}
}

func (c *Clockwerk) runPending() {
	for _, entry := range c.entries {
		if time.Now().Before(entry.Next) {
			break
		}

		c.schedule(entry)
		go c.runJob(entry)
	}
}

func (c *Clockwerk) runJob(e *Entry) {
	defer func() {
		// eat the panic
		if r := recover(); r != nil {
		}
	}()

	e.Job.Run(c.ctx)
}
