// internal/app/services/background_job_service/background_job.go
package background_job_service

import (
	"clean_architecture_go/internal/app/services"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type JobQueue struct {
	jobs          chan services.Job
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        context.CancelFunc
	maxWorkers    int
	activeWorkers atomic.Int32
	idleTimeout   time.Duration
}

func NewJobQueue(maxWorkers int, queueSize int) *JobQueue {
	ctx, cancel := context.WithCancel(context.Background())

	job := &JobQueue{
		jobs:        make(chan services.Job, queueSize),
		ctx:         ctx,
		cancel:      cancel,
		maxWorkers:  maxWorkers,
		idleTimeout: 30 * time.Second, // Kill idle workers after 30s
	}

	job.Start()
	return job
}

func (q *JobQueue) Start() {
	// Start with 1 worker minimum
	q.spawnWorker()
	log.Printf("✅ Auto-scaling job queue started (max %d workers)", q.maxWorkers)
}

func (q *JobQueue) spawnWorker() {
	if int(q.activeWorkers.Load()) >= q.maxWorkers {
		return // Max workers reached
	}

	q.activeWorkers.Add(1)
	q.wg.Add(1)
	go q.worker()
}

func (q *JobQueue) worker() {
	defer q.wg.Done()
	defer q.activeWorkers.Add(-1)

	log.Printf("👷 Worker started (active: %d)", q.activeWorkers.Load())

	for {
		select {
		case <-q.ctx.Done():
			return

		case job, ok := <-q.jobs:
			if !ok {
				return
			}
			if err := job(q.ctx); err != nil {
				log.Printf("🔴 Job failed: %v", err)
			}

		case <-time.After(q.idleTimeout):
			// Worker idle for too long, die if not the last one
			if q.activeWorkers.Load() > 1 {
				log.Printf("👷 Worker exiting (idle, active: %d)", q.activeWorkers.Load()-1)
				return
			}
		}
	}
}

func (q *JobQueue) Enqueue(job services.Job) {
	// If queue has pending jobs, try to spawn more workers
	if len(q.jobs) > 0 && int(q.activeWorkers.Load()) < q.maxWorkers {
		q.spawnWorker()
	}
	log.Printf("📥 Enqueue called, pending before: %d", len(q.jobs))
	
	q.jobs <- job
	
	log.Printf("📥 Enqueue done, pending after: %d", len(q.jobs))
}

func (q *JobQueue) EnqueueWithTimeout(job services.Job) bool {
	select {
	case q.jobs <- job:
		return true
	default:
		return false
	}
}

func (q *JobQueue) Shutdown() {
	q.cancel()
	close(q.jobs)
	q.wg.Wait()
	log.Println("✅ Job queue shutdown complete")
}

func (q *JobQueue) PendingJobs() int {
	return len(q.jobs)
}

func (q *JobQueue) ActiveWorkers() int {
	return int(q.activeWorkers.Load())
}
