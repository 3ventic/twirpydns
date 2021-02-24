package workers

import "sync"

type Worker interface {
	Run(id string, fn func() interface{}) (<-chan interface{}, bool)
}

type worker struct {
	mu   sync.Mutex
	jobs map[string]*job
}

func New() Worker {
	return &worker{
		mu:   sync.Mutex{},
		jobs: make(map[string]*job),
	}
}

// Run starts or subscribes to an existing task based on the given id.
func (w *worker) Run(id string, fn func() interface{}) (<-chan interface{}, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	j, ok := w.jobs[id]
	if ok {
		j.subscribers++
		return j.channel, false
	}

	c := make(chan interface{})
	w.jobs[id] = &job{
		channel:     c,
		subscribers: 1,
	}

	go func() {
		result := fn()

		w.mu.Lock()
		defer w.mu.Unlock()

		w.jobs[id].BroadcastResult(result)
		delete(w.jobs, id)
	}()

	return c, true
}
