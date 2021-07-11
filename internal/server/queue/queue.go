package queue

import (
	"time"

	"github.com/groupe-edf/watchdog/internal/logging"
	"github.com/groupe-edf/watchdog/internal/models"
)

type Client struct {
	Logger logging.Interface
	store  models.Store
}

func (client *Client) Enqueue(job *models.Job) error {
	if job.StartedAt.IsZero() {
		job.StartedAt = time.Now()
	}
	return client.store.Enqueue(job)
}

func (client *Client) LockJob(queue string) (*models.Job, error) {
	return client.store.LockJob(queue)
}

func NewClient(store models.Store, logger logging.Interface) *Client {
	return &Client{
		Logger: logger,
		store:  store,
	}
}
