// SPDX-FileCopyrightText: Copyright (c) 2024, Obviously, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package workerpool_test

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-obvious/workerpool"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("when poolSize is less than 1", func(t *testing.T) {
		inst := workerpool.New(0, 10)
		assert.Nil(t, inst)
	})

	t.Run("when bufferSize is negative", func(t *testing.T) {
		inst := workerpool.New(1, -1)
		assert.Nil(t, inst)
	})

	t.Run("when poolSize and bufferSize are valid", func(t *testing.T) {
		inst := workerpool.New(2, 10)
		assert.NotNil(t, inst)
		inst.Stop()
	})
}

func TestStartAndStop(t *testing.T) {
	t.Parallel()

	inst := workerpool.New(2, 10)
	assert.NotNil(t, inst)

	var wg sync.WaitGroup
	wg.Add(1)

	inst.Start()

	// add some work to the pool
	inst.WorkCh <- func() {
		time.Sleep(100 * time.Millisecond)
		wg.Done()
	}

	// wait for the work to be done
	wg.Wait()

	inst.Stop()
}

// These are good examples to show how to leverage the workerpool
func TestExample1(t *testing.T) {
	t.Parallel()
	const itemCount = 1000

	inst := workerpool.New(10, 10)
	assert.NotNil(t, inst)
	//
	inst.Start()
	defer inst.Stop()
	//
	var wg sync.WaitGroup
	values := make([]int, itemCount+1) // add one more for safety
	for i := 0; i < itemCount; i++ {
		// This MUST BE here to get a local stack value which only the
		// routing below owns. It can't be in the method def
		i := i
		pipevalue := func() {
			defer wg.Done()
			values[i] = i
		}
		wg.Add(1)
		inst.WorkCh <- pipevalue
	}
	wg.Wait()
	//
	// now validate
	for i := 0; i < itemCount; i++ {
		assert.Equal(t, i, values[i])
	}
}

func TestExample2(t *testing.T) {
	t.Parallel()
	inst := workerpool.New(10, 10)
	assert.NotNil(t, inst)
	inst.Start()
	defer inst.Stop()
	var wg sync.WaitGroup
	ch := make(chan int, 1000)
	done := make(chan bool)

	for i := 0; i < 1000; i++ {
		// This MUST BE here to get a local stack value which only the
		// routing below owns. It can't be in the method def
		i := i
		pipevalue := func() {
			// i := i // NOT HERE
			defer wg.Done()
			ch <- i
		}
		wg.Add(1)
		inst.WorkCh <- pipevalue
	}

	got := 0
	go func() {
		for val := range ch {
			if val < 0 {
				break
			}
			got++
		}
		done <- true
	}()

	// wait for the work to be done
	wg.Wait()
	ch <- -1
	<-done
	assert.Equal(t, 1000, got)
}

// UNTESTABLE
// func TestStartRecover(t *testing.T) {
// 	t.Parallel()
// 	panicTest := func() { // create a new instance with pool size 1 and buffer size 0
// 		inst := workerpool.New(1, 0)
// 		assert.NotNil(t, inst)

// 		// start the instance
// 		inst.Start()

// 		// add a panic-inducing work function to the pool
// 		inst.WorkCh <- func() {
// 			panic("something went wrong")
// 		}

// 		// wait for the panic to be handled by the recover code
// 		time.Sleep(100 * time.Millisecond)

// 		// stop the instance
// 		inst.Stop()
// 	}
// 	assert.Panics(t, panicTest, "The code did not panic")
// }
