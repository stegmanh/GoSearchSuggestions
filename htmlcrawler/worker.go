package htmlcrawler

import (
	"fmt"
)

type Worker struct {
	Work        chan string
	WorkerQueue chan chan string
	Commands    chan string
}

func CreateWorker(workerQueue chan chan string) Worker {
	worker := Worker{
		Work:        make(chan string),
		WorkerQueue: workerQueue,
		Commands:    make(chan string),
	}

	return worker
}

func (this Worker) Start() {
	go func() {
		for {
			this.WorkerQueue <- this.Work

			select {
			case work := <-this.Work:
				fmt.Println(work)
				//Crawl the urls
			case command := <-this.Commands:
				switch command {
				case "stop":
					fmt.Println(command)
				case "start":
					fmt.Println(command)
				default:
					fmt.Println(command)
				}
			}
		}
	}()
}

func (this Worker) Stop() {
	go func() {
		this.Commands <- "stop"
	}()
}
