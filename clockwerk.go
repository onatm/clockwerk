package clockwerk

import (
	"time"
)

type Job interface {
	Run()
}

type Entry struct {
	Period time.Duration
	Next   time.Time
	Prev   time.Time
	Job    Job
}

func newEntry(period time.Duration) *Entry {
	return &Entry{
		Period: period,
		Prev:   time.Unix(0, 0),
	}
}

func (e *Entry) Do(job Job) {
	e.Job = job
}

type Clockwerk struct {
	entries []*Entry
	stop    chan struct{}
	running bool
}

func New() *Clockwerk {
	return &Clockwerk{
		entries: nil,
		stop:    make(chan struct{}),
		running: false,
	}
}

func (c *Clockwerk) EverySeconds(seconds uint64) *Entry {
	entry := newEntry(time.Duration(seconds) * time.Second)

	c.schedule(entry)

	if !c.running {
		c.entries = append(c.entries, entry)
	}

	return entry
}

func (c *Clockwerk) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Clockwerk) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func (c *Clockwerk) schedule(e *Entry) {
	e.Prev = time.Now()

	e.Next = e.Prev.Add(e.Period)
}

func (c *Clockwerk) run() {
	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				c.runPending()
				continue
			case <-c.stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (c *Clockwerk) runPending() {
	go func() {
		for _, entry := range c.entries {
			if time.Now().After(entry.Next) {
				go c.runJob(entry)
			}
		}
	}()
}

func (c *Clockwerk) runJob(e *Entry) {
	defer func() {
		if r := recover(); r != nil {
			c.schedule(e)
		}
	}()

	c.schedule(e)
	e.Job.Run()
}
