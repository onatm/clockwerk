/*
Package clockwerk implements an in-process scheduler for periodic jobs.

Usage

Callers may register Jobs to be invoked on a given schedule. Clockwerk will run
them in their own goroutines.

  type DummyJob struct{}

  func (d DummyJob) Run() {
    fmt.Println("Every 30 seconds")
  }
  ...
  var job DummyJob
  c := clockwerk.New()
  c.Every(30 * time.Second).Do(job)
  c.Start()
  ...
  // Funcs are invoked in their own goroutine, asynchronously.
  ...
  c.Stop()  // Stop the scheduler (does not stop any jobs already running).

*/
package clockwerk
