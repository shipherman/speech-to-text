package cmd

type Queue struct {
	C chan struct{}
}

func (s *Queue) Acquire() {
	s.C <- struct{}{}
}

func (s *Queue) Release() {
	<-s.C
}
