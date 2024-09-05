package semaphore

type Semaphore struct {
	SemC chan struct{}
}

func (s *Semaphore) Acquire() {
	s.SemC <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.SemC
}
