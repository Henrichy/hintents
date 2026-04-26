// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

package watch

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	frames    []string
	current   int
	done      chan struct{}
	mu        sync.Mutex
	isRunning bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"|", "/", "-", "\\"},
		done:   make(chan struct{}),
	}
}

func (s *Spinner) Start(message string) {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return
	}
	s.isRunning = true
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				fmt.Print("\r\033[K")
				return
			case <-ticker.C:
				s.mu.Lock()
				fmt.Printf("\r%s %s", s.frames[s.current], message)
				s.current = (s.current + 1) % len(s.frames)
				s.mu.Unlock()
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.isRunning {
		s.mu.Unlock()
		return
	}
	s.isRunning = false
	s.mu.Unlock()

	// Signal the goroutine to exit.
	select {
	case s.done <- struct{}{}:
	default:
		// If the channel is full, someone already signaled it
		// or the goroutine already exited.
	}

	// Give the goroutine a moment to receive and clear the line.
	// In a real implementation, we might want a sync.WaitGroup here.
	time.Sleep(20 * time.Millisecond)
}

func (s *Spinner) StopWithMessage(message string) {
	s.Stop()
	fmt.Printf("\r[OK] %s\n", message)
}

func (s *Spinner) StopWithError(message string) {
	s.Stop()
	fmt.Printf("\r[ERROR] %s\n", message)
}
