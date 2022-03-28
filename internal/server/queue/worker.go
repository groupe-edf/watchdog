package queue

import (
	"sync"
	"time"

	"github.com/groupe-edf/watchdog/internal/core/models"
	"github.com/groupe-edf/watchdog/pkg/logging"
)

const (
	MaxErroCount = 3
)

type Worker struct {
	Logger      logging.Interface
	Queue       string
	WorkMap     WorkMap
	channel     chan struct{}
	client      *Client
	done        bool
	idleTimeout time.Duration
	lock        sync.Mutex
	jobs        []models.Job
}

func (worker *Worker) Work() {
	for {
		if worker.WorkOne() {
			select {
			case <-worker.channel:
				return
			default:
			}
		} else {
			select {
			case <-worker.channel:
				return
			case <-time.After(worker.idleTimeout):
			}
		}
	}
}

func (worker *Worker) WorkOne() (done bool) {
	job, err := worker.client.store.LockJob(worker.Queue)
	if err != nil {
		return
	}
	if job == nil {
		return
	}
	defer worker.client.store.DoneJob(job)
	defer func() {
		if r := recover(); r != nil {
			if err := worker.client.store.SaveJobError(job, "recovering error"); err != nil {
				worker.Logger.Error(err)
			}
		}
	}()
	if job.ErrorCount < 3 {
		workFunc, ok := worker.WorkMap[job.Type]
		if !ok {
			return
		}
		if err = workFunc(job); err != nil {
			worker.Logger.Error(err)
			worker.client.store.SaveJobError(job, err.Error())
			return
		}
	}
	if err = worker.client.store.DeleteJob(job); err != nil {
		worker.Logger.Infof("attempting to delete job %d: %v", job.ID, err)
	}
	return
}

func (worker *Worker) Shutdown() {
	worker.lock.Lock()
	defer worker.lock.Unlock()
	if worker.done {
		return
	}
	worker.channel <- struct{}{}
	worker.done = true
	close(worker.channel)
}

func NewWorker(client *Client, workMap WorkMap) *Worker {
	return &Worker{
		Queue:       "default",
		WorkMap:     workMap,
		client:      client,
		idleTimeout: 5 * time.Second,
	}
}
