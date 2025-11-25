/**
* @license
* Copyright 2025 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

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
