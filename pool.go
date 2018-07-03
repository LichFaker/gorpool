package gorpool

import "sync"

const (
    defaultWorkerCount = 100 // the default count of the worker
)

var (
    oncePool    sync.Once // perform exactly one action
    defaultPool *Pool     // the singleton object of Pool
)

// Task
type Task func()

// Pool
type Pool struct {
    // the count of the worker need to handle the task.
    // default value is 100.
    workerCount uint32
    
    // the channel for worker
    workerChan chan chan Task
    // the channel for task
    taskChan chan Task
    // the exit single
    exitChan chan struct{}
}

// GetPool get the Pool's singleton
func GetPool() *Pool {
    if defaultPool == nil {
        oncePool.Do(func() {
            defaultPool = NewPool()
        })
    }
    return defaultPool
}

// NewPool create and return the pointer of Pool
func NewPool() *Pool {
    pool := new(Pool)
    pool.init()
    return pool
}

// initialization
func (p *Pool) init() {
    // set the default count of worker
    p.workerCount = defaultWorkerCount
    // allocates and initializes the channel
    p.workerChan = make(chan chan Task)
    p.taskChan = make(chan Task)
    p.exitChan = make(chan struct{})
}

// Start handle the task and do it in an appropriate worker.
// Create some workers which to do the task in goroutine.
// Once the task is received, just put it into the task slice,
// so that the current goroutine won't block.
// Only if the slice of task is not empty and has active worker,
// then the task in the head of task slice will be send to the appropriate worker.
func (p *Pool) Start() {
    go func() {
        // create workers
        for i := uint32(0); i < p.workerCount; i++ {
            p.createWorker()
        }
        
        // the queue to store the active worker
        // the value type is chan Task
        workerQ := Queue{}
        // the queue to store the task
        // the value type is Task
        taskQ := Queue{}
        
        var (
            activeWorker chan Task
            activeTask   Task
        )
    
    Loop:
        for {
            if !workerQ.IsEmpty() && !taskQ.IsEmpty() {
                activeWorker = workerQ.GetTop().(chan Task)
                activeTask = taskQ.GetTop().(Task)
            } else {
                activeTask = nil
                activeWorker = nil
            }
            
            select {
            case task := <-p.taskChan:
                taskQ.Push(task)
            case worker := <-p.workerChan:
                workerQ.Push(worker)
            case activeWorker <- activeTask:
                workerQ.Pop()
                taskQ.Pop()
            case <-p.exitChan:
                break Loop
            }
        }
    }()
}

// SetWorkerCount set the count of the worker.
func (p *Pool) SetWorkerCount(workerCount uint32) {
    p.workerCount = workerCount
}

// Submit submit a task
func (p *Pool) Submit(task Task) {
    p.taskChan <- task
}

// Destroy destroy and exit the pool
func (p *Pool) Destroy() {
    var sig struct{}
    p.exitChan <- sig
}

// create the worker to do the task
func (p *Pool) createWorker() {
    in := make(chan Task)
    go func() {
        for {
            p.workerNotify(in)
            task := <-in
            task()
        }
    }()
}

// notify the worker is active
func (p *Pool) workerNotify(wc chan Task) {
    p.workerChan <- wc
}
