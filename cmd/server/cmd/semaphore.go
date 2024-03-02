package cmd

import (
	"time"

	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
)

type Queue struct {
	C chan struct{}
}

type resultWithError struct {
	A   sttservice.Audio
	Err error
}

func (s *Queue) Acquire() {
	s.C <- struct{}{}
}

func (s *Queue) Release() {
	time.Sleep(10 * time.Second)
	<-s.C
}
