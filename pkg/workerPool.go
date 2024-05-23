package pkg

type WorkerPool struct {
	pool   chan struct{}
	worker chan func()
}

func NewWorkerPool(poolSize int) *WorkerPool {
	return &WorkerPool{
		pool:   make(chan struct{}, poolSize),
		worker: make(chan func()),
	}
}

func (wp *WorkerPool) AddWorker(fn func()) {

	wp.pool <- struct{}{}

	go func() {

		worker := <-wp.worker

		worker()
		<-wp.pool

	}()

	wp.worker <- fn

}
