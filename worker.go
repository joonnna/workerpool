package workerpool

type worker struct {
	workChan chan *work
	poolChan chan chan *work
}

func newWorker(poolChan chan chan *work) *worker {
	return &worker{
		workChan: make(chan *work),
		poolChan: poolChan,
	}
}

func (w *worker) start() {
	go func() {
		for {
			w.poolChan <- w.workChan

			select {
			case w, ok := <-w.workChan:
				if !ok {
					return
				}
				w.payload()
			}
		}
	}()
}
