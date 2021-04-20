package download

type Manager struct {
	workQueues []workQueue
	results    resultQueue
}

type workQueue chan interface{}
type resultQueue chan interface{}

type Config struct {
	WorkQueueSize  int
	TotalWorkQueue int
}

func NewManager(cfg Config) *Manager {
	workQueues := make([]workQueue, cfg.TotalWorkQueue)
	for i := range workQueues {
		workQueues[i] = make(workQueue, cfg.WorkQueueSize)
	}

	results := make(resultQueue)
	return &Manager{
		workQueues,
		results,
	}
}
