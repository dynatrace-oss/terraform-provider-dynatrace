package shutdown

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var System = newone()

func newone() *system {
	syst := &system{running: true}
	syst.Wait()
	return syst
}

type system struct {
	mu      sync.Mutex
	running bool
}

func (me *system) Running() bool {
	me.mu.Lock()
	defer me.mu.Unlock()
	return me.running
}

func (me *system) Stopped() bool {
	me.mu.Lock()
	defer me.mu.Unlock()
	return !me.running
}

func (me *system) Stop() {
	me.mu.Lock()
	defer me.mu.Unlock()
	me.running = false
}

func (me *system) handler(signal os.Signal) {
	if signal == syscall.SIGINT {
		me.Stop()
	}
}

func (me *system) Wait() {
	sigchn := make(chan os.Signal, 1)
	signal.Notify(sigchn)

	go func() {
		for {
			s := <-sigchn
			me.handler(s)
		}
	}()
}
