// SPDX-FileCopyrightText: Copyright (c) 2016-2024, CloudZero, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package workerpool

import (
	"fmt"
	"os"
	"runtime/debug"
)

type Instance struct {
	poolSize int
	done     chan struct{}
	WorkCh   chan func()
}

func New(poolSize, bufferSize int) *Instance {
	if poolSize < 1 {
		return nil
	}

	if bufferSize < 0 {
		return nil
	}

	inst := &Instance{
		poolSize: poolSize,
		done:     make(chan struct{}),
		WorkCh:   make(chan func(), bufferSize),
	}

	return inst
}

func (ref *Instance) Start() {
	for i := 0; i < ref.poolSize; i++ {
		go func(id int) {
			// log panics
			defer func() {
				if r := recover(); r != nil {
					stack := string(debug.Stack())
					fmt.Println("panic:", r, "\n"+stack)
					os.Exit(1)
				}
			}()
			for {
				select {
				case <-ref.done:
					return
				case work, ok := <-ref.WorkCh:
					if ok && work != nil {
						work()
					}
				}
			}
		}(i)
	}
}

func (ref *Instance) Stop() {
	close(ref.done)
}
