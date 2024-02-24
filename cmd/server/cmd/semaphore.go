package cmd

import (
	sttservice "github.com/shipherman/speech-to-text/gen/stt/service/v1"
)

type Semaphore struct {
	C chan struct{}
}

type resultWithError struct {
	A   sttservice.Audio
	Err error
}

func (s *Semaphore) Acquire() {
	s.C <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.C
}

// func ProcessAudio(a []*sttservice.Audio, pool int) error {
// 	// sem := Semaphore{
// 	// 	C: make(chan struct{}, pool),
// 	// }

// 	// outputCh := make(chan resultWithError, len(a))

// 	// sgnlCh := make(chan struct{})

// 	// output := make([]*sttservice.Audio, 0, len(a))

// 	return nil
// }
