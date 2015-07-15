package htmlcrawler

import ()

var WorkerQueue chan chan string

func StartWorkers(numWorkers int) {
	WorkerQueue = make(chan chan string, numWorkers)

	for i := 0; i < numWorkers; i++ {
		worker := CreateWorker(workerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
				case: work := <-
			}
		}
	}
}
