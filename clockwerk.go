package clockwerk

import (
	"time"
)

// Job is an interface for added jobs.
type Job interface {
	Run()
}

// Entry consists of a schedule and the job to execute on that schedule.
type Entry struct {
	Period time.Duration
	Next   time.Time
	Job    Job
}

func newEntry(period time.Duration) *Entry {
	return &Entry{
		Period: period,
	}
}

// Do adds a Job to the Entry.
func (e *Entry) Do(job Job) {
	e.Job = job
}

// Clockwerk keeps track of any number of entries, invoking associated Job's Run
// method as specified by the schedule.
type Clockwerk struct {
	entries []*Entry
	stop    chan struct{}
	running bool
}

// New returns a new Clockwerk job runner.
func New() *Clockwerk {
	return &Clockwerk{
		entries: nil,
		stop:    make(chan struct{}),
		running: false,
	}
}

// Every schedules a new Entry and returns it.
func (c *Clockwerk) Every(period time.Duration) *Entry {
	entry := newEntry(period)

	c.schedule(entry)

	if !c.running {
		c.entries = append(c.entries, entry)
	}

	return entry
}

// Start the Clockwerk in its own go-routine, or no-op if already started.
func (c *Clockwerk) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

// Stop the Clockwerk if it is running.
func (c *Clockwerk) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func (c *Clockwerk) schedule(e *Entry) {
	e.Next = time.Now().Add(e.Period)
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
