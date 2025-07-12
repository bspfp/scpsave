package sem

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(n int) *Semaphore {
	if n <= 0 {
		n = 1
	}
	return &Semaphore{
		ch: make(chan struct{}, n),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) Release() {
	select {
	case <-s.ch:
	default:
		panic("Semaphore release called without a corresponding acquire")
	}
}
