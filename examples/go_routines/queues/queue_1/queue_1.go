/*
Queue 1
Very simple example showcasing how an HTTP server could handle the requests for learning purposes.
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	WorkersCount            = 3
	MaxRequestIdlePerWorker = 100
	Debug                   = false
)

type Job func(reqId int) string

type Request struct {
	ID       int
	Response chan string
}

type Worker struct {
	Name   string
	work   <-chan *Request
	logger *log.Logger
}

func (w *Worker) Work(wg *sync.WaitGroup, job Job) {
	defer wg.Done()
	w.logger.Println(" Started, waiting for request ...")

	for {
		select {
		case req, ok := <-w.work:
			if !ok {
				w.logger.Println("Request pool close, exiting ...")
				return
			}
			w.handleRequest(req, job)

		default:
			log.Println("Worker is idle, waiting for request ...")
			time.Sleep(120 * time.Millisecond)
		}
	}
}

func (w *Worker) handleRequest(req *Request, job Job) {
	if Debug {
		w.logger.Printf("Received request %d\n", req.ID)
	}
	response := job(req.ID)
	req.Response <- response
	close(req.Response)
}

type RequestPool struct {
	ctx     context.Context
	cancel  context.CancelFunc
	job     Job
	pool    chan *Request
	workers []*Worker
}

func NewRequestPool(workers int, job Job) *RequestPool {
	// NOTE:
	// We need to have a HIGHER buffer in the request pool in order to
	// avoid losing requests | effectively this is the "Max request waiting per worker" which is something commonly seen
	// in a web server. we have I workers and J max idle requests.
	pool := make(chan *Request, MaxRequestIdlePerWorker)

	w := make([]*Worker, workers)
	for i := range workers {
		name := fmt.Sprintf("worker_%d", i)
		w[i] = &Worker{
			Name:   name,
			work:   pool,
			logger: log.New(os.Stderr, name+" ", 0),
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &RequestPool{
		ctx:     ctx,
		cancel:  cancel,
		workers: w,
		job:     job,
		pool:    pool,
	}
}

func (p *RequestPool) StartWorkers(wg *sync.WaitGroup) {
	for _, w := range p.workers {
		wg.Add(1)
		go w.Work(wg, p.job)
	}
}

func (p *RequestPool) Shutdown() {
	log.Println("Shutting down RequestPool...")

	p.cancel() // That's what I meant when closing by using a ctx... not sure if is the best way for this problem...0
	time.Sleep(10 * time.Millisecond)

	// I understand this now... this serves as a "signal" for the worker to know that "no more request are accepted"
	close(p.pool)

}

func (p *RequestPool) SendSyncRequest(reqId int) (string, error) {
	request := &Request{
		ID:       reqId,
		Response: make(chan string),
	}

	select {
	case <-p.ctx.Done():
		log.Println("RequestPool is shutting down, Dropping new requests")
		return "", errors.New("server Shutting Down")
	case p.pool <- request:

	}

	response := <-request.Response
	return response, nil
}

func job(reqId int) string {
	result := fmt.Sprintf("!!!!RESPONE (REQ %d)!!!", reqId)
	time.Sleep(120 * time.Millisecond)
	return result
}

func main() {
	// Main uses default logger
	log.SetFlags(0)
	log.SetPrefix("main: ")
	log.Println("Starting ...")

	var wg sync.WaitGroup

	requestPool := NewRequestPool(WorkersCount, job)
	requestPool.StartWorkers(&wg)

	totalClients := 50
	for i := range totalClients {
		go func() {
			if Debug {
				log.Printf("Client %d -- Sending Request ....", i)
			}
			response, err := requestPool.SendSyncRequest(i)
			if err != nil {
				log.Printf("Client %d -- Error: %s", i, err.Error())
				return
			}

			log.Printf("Client %d -- Response: %s", i, response)
		}()
	}

	timeToWait := 500 * time.Millisecond
	log.Printf(
		"Sent %d requests, waiting %d Give some room for workers to process requests as \"normal\"",
		totalClients,
		timeToWait,
	)
	time.Sleep(timeToWait)

	requestPool.Shutdown()

	log.Println("Waiting for all workers to finish...")
	wg.Wait()
	log.Println("All workers finished!")

}
