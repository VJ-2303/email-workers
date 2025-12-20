package worker

import (
	"log"
	"sync"
)

type Task struct {
	From    string
	To      string
	Subject string
	Body    string
}

type Pool struct {
	maxWorkers int
	taskQuese  chan Task
	wg         sync.WaitGroup
	logger     *log.Logger
}

func NewPool(maxWorkers int, queueSize int, logger *log.Logger) *Pool {
	return &Pool{
		maxWorkers: maxWorkers,
		taskQuese:  make(chan Task, queueSize),
		logger:     logger,
	}
}

func (p *Pool) Run() {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i + 1)
	}
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	p.logger.Printf("Worker %d started", id)

	for task := range p.taskQuese {
		p.logger.Printf("Worker %d processing email to %s", id, task.To)

		p.process(task)
	}
	p.logger.Printf("Worker %d stopped", id)
}

func (p *Pool) process(t Task) {
	p.logger.Printf("Sending email: Subject=%q", t.Subject)
}

func (p *Pool) Submit(t Task) {
	p.taskQuese <- t
}
