package pool

import (
	"fmt"
	"sync"

	"github.com/mlcdf/sc-backup/internal/sc"
)

type RunFunc func() (interface{}, error)

// Task encapsulates a work item that should go in a work
// pool.
type Task struct {
	// Err holds an error that occurred during a task. Its
	// result is only meaningful after Run has been called
	// for the pool that holds it.
	Err  error
	Out  interface{}
	Func RunFunc
}

// NewTask initializes a new task based on a given work
// function.
func NewTask(f RunFunc) *Task {
	return &Task{Func: f}
}

// Run runs a Task and does appropriate accounting via a
// given sync.WorkGroup.
func (t *Task) Run(wg *sync.WaitGroup) {
	t.Out, t.Err = t.Func()
	wg.Done()
}

// Pool is a worker group that runs a number of tasks at a
// configured concurrency.
type Pool struct {
	Tasks []*Task

	concurrency int
	tasksChan   chan *Task
	wg          sync.WaitGroup
}

// NewPool initializes a new pool with the given tasks and
// at the given concurrency.
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		tasksChan:   make(chan *Task),
	}
}

// Run runs all work within the pool and blocks until it's
// finished.
func (p *Pool) Run() {
	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
}

// The work loop for any single goroutine.
func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}

// Merge the tasks result
func (p *Pool) Merge(entries []*sc.Entry) ([]*sc.Entry, error) {
	for _, task := range p.Tasks {
		if task.Err != nil {
			return nil, task.Err
		}

		_out, ok := task.Out.([]*sc.Entry)
		if !ok {
			return nil, fmt.Errorf("critical: failed to cast to []*Entry. Please open a bug report at https://github.com/mlcdf/sc-backup")
		}
		entries = append(entries, _out...)
	}
	return entries, nil
}
