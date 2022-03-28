package queue

import (
	"sync"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

const (
	defaultQueue = "default"
)

type Options struct {
	MaxWorkers   int
	MaxQueueSize int
}

type WorkerPool struct {
	Logger  logging.Interface
	WorkMap WorkMap
	client  *Client
	done    bool
	lock    sync.Mutex
	queue   string
	workers []*Worker
}

type WorkFunc func(job *models.Job) error

type WorkMap map[string]WorkFunc

func (pool *WorkerPool) Enqueue(job *models.Job) error {
	return pool.client.Enqueue(job)
}

func (pool *WorkerPool) Start() {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	for i := range pool.workers {
		pool.workers[i] = NewWorker(pool.client, pool.WorkMap)
		pool.workers[i].Queue = pool.queue
		pool.workers[i].Logger = pool.client.Logger
		go pool.workers[i].Work()
	}
}

func (pool *WorkerPool) Shutdown() {
	pool.lock.Lock()
	defer pool.lock.Unlock()
	if pool.done {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(pool.workers))
	for _, worker := range pool.workers {
		go func(worker *Worker) {
			if worker != nil {
				worker.Shutdown()
			}
			wg.Done()
		}(worker)
	}
	wg.Wait()
	pool.done = true
}

func NewWorkerPool(client *Client, workMap WorkMap, options Options) *WorkerPool {
	return &WorkerPool{
		queue:   defaultQueue,
		client:  client,
		workers: make([]*Worker, options.MaxWorkers),
		WorkMap: workMap,
	}
}
