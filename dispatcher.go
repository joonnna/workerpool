package workerpool

type work struct {
	payload func()
}

type Dispatcher struct {
	exitChan   chan bool
	jobChan    chan *work
	maxWorkers uint32

	pool chan chan *work
}

func NewDispatcher(maxWorkers uint32) *Dispatcher {
	return &Dispatcher{
		maxWorkers: maxWorkers,
		pool:       make(chan chan *work, maxWorkers),
		jobChan:    make(chan *work),
		exitChan:   make(chan bool),
	}
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case j := <-d.jobChan:
			go func() {
				workerChan := <-d.pool
				workerChan <- j
			}()
		case <-d.exitChan:
			return
		}
	}
}

func (d *Dispatcher) Submit(job func()) {
	d.jobChan <- &work{payload: job}
}

func (d *Dispatcher) Start() {
	var i uint32
	for i = 0; i < d.maxWorkers; i++ {
		w := newWorker(d.pool)
		w.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) Stop() {
	go func() {
		close(d.exitChan)

		var i uint32 = 0

		for i = 0; i < d.maxWorkers; i++ {
			j := <-d.pool
			close(j)
		}
	}()
}
