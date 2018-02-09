package workerpool

type worker struct {
	exitChan chan bool
	workChan chan *work
	poolChan chan chan *work
}

func newWorker(poolChan chan chan *work, exitChan chan bool) *worker {
	return &worker{
		workChan: make(chan *work),
		exitChan: exitChan,
		poolChan: poolChan,
	}
}

func (w *worker) start() {
	go func() {
		for {
			w.poolChan <- w.workChan

			select {
			case w := <-w.workChan:
				w.payload()

			case <-w.exitChan:
				return
			}
		}
	}()
}
