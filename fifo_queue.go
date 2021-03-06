package goconcurrentqueue

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

// FIFO (First In First Out) concurrent queue
type FIFO struct {
	slice       []interface{}
	rwmutex     sync.RWMutex
	lockRWmutex sync.RWMutex
	isLocked    bool
}

// NewFIFO returns a new FIFO concurrent queue
func NewFIFO() *FIFO {
	ret := &FIFO{}
	ret.initialize()

	return ret
}

func (st *FIFO) initialize() {
	st.slice = make([]interface{}, 0)
}

// Enqueue enqueues an element
func (st *FIFO) Enqueue(value interface{}) error {
	if st.isLocked {
		return errors.New("The queue is locked")
	}

	st.rwmutex.Lock()
	defer st.rwmutex.Unlock()

	st.slice = append(st.slice, value)
	return nil
}

// Dequeue dequeues an element
func (st *FIFO) Dequeue() (interface{}, error) {
	if st.isLocked {
		return nil, errors.New("The queue is locked")
	}

	st.rwmutex.Lock()
	defer st.rwmutex.Unlock()

	len := len(st.slice)
	if len == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	elementToReturn := st.slice[0]
	st.slice = st.slice[1:]

	return elementToReturn, nil
}

// Get returns an element's value and keeps the element at the queue
func (st *FIFO) Get(index int) (interface{}, error) {
	if st.isLocked {
		return nil, errors.New("The queue is locked")
	}

	st.rwmutex.RLock()
	defer st.rwmutex.RUnlock()

	if len(st.slice) <= index {
		return nil, fmt.Errorf("index out of bounds: %v", index)
	}

	return st.slice[index], nil
}

// Remove removes an element from the queue
func (st *FIFO) Remove(index int) error {
	if st.isLocked {
		return errors.New("The queue is locked")
	}

	st.rwmutex.Lock()
	defer st.rwmutex.Unlock()

	if len(st.slice) <= index {
		return fmt.Errorf("index out of bounds: %v", index)
	}

	// remove the element
	st.slice = append(st.slice[:index], st.slice[index+1:]...)

	return nil
}

// GetLen returns the number of enqueued elements
func (st *FIFO) GetLen() int {
	st.rwmutex.RLock()
	defer st.rwmutex.RUnlock()

	return len(st.slice)
}

// Lock // Locks the queue. No enqueue/dequeue operations will be allowed after this point.
func (st *FIFO) Lock() {
	st.lockRWmutex.Lock()
	defer st.lockRWmutex.Unlock()

	st.isLocked = true
}

// Unlock unlocks the queue
func (st *FIFO) Unlock() {
	st.lockRWmutex.Lock()
	defer st.lockRWmutex.Unlock()

	st.isLocked = false
}

// IsLocked returns true whether the queue is locked
func (st *FIFO) IsLocked() bool {
	st.lockRWmutex.RLock()
	defer st.lockRWmutex.RUnlock()

	return st.isLocked
}
