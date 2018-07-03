package gorpool

import (
	"context"
)

type WorkerNotifier interface {
	Notify(w *Worker)
}

type Worker struct {
	in       chan Task
	notifier WorkerNotifier
}

// create the worker to do the task
func NewWorker(notifier WorkerNotifier) *Worker {
	w := new(Worker)
	w.init()
	w.notifier = notifier
	return w
}

func (w *Worker) init() {
	w.in = make(chan Task)
}

func (w *Worker) GetChan() chan Task {
	return w.in
}

func (w *Worker) Run(ctx context.Context) {
	go func() {
		for {
			w.notifier.Notify(w)
			select {
			case <-ctx.Done():
				return
			case task := <-w.in:
				task.Do()
			}
			
		}
	}()
}
