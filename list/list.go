// Copyright 2014 Pan Qing. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package list implements a effective, multiple kinds and goroutine safe list.
//
package list

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// A Kind represents the specific kind of list.
type Kind uint

const (
	Invalid Kind = iota
	FIFO         // First in, last out
	LIFO         // First in, first out
	RAND         // Random
)

var KindNames = []string{
	Invalid: "Invalid",
	FIFO:    "FIFO",
	LIFO:    "LIFO",
	RAND:    "RANDOM",
}

func (k Kind) String() string {
	if int(k) < len(KindNames) {
		return KindNames[k]
	}
	return "kind" + strconv.Itoa(int(k))
}

// Bucket
type bucket struct {
	size         int   // Elements number in bucket
	popPosition  int   // First elmenent position for pop. If popPosition equals -1 then bucket is empty.
	pushPosition int   // Last element position + 1 for push. If pushPosition equals size then bucket is full.
	elems        Elems // An elements slice with bucketSize capbility
}

func newBucket(size int, e Elem) (b *bucket) {
	b = new(bucket)
	b.size = size
	b.popPosition = -1 // -1 means no element for p
	b.elems = e.Array(size)
	return
}

// List represents a FIFO|LIFO|RAND list.
type List struct {
	kind              Kind      // list kind
	count             int       // current list length
	bucketSize        int       // Elements number in bucket, it is fixed after initial.
	currentPushBucket int       // Which bucket for push
	pool              []*bucket // buckets
	r                 *rand.Rand
	mu                sync.Mutex
}

// New returns an initialized list.
func New(kind Kind, bucketSize int) (l *List, err error) {
	l = new(List)
	l.mu.Lock()
	defer l.mu.Unlock()

	switch kind {
	case FIFO:
		fallthrough
	case LIFO:
		l.kind = kind
	case RAND:
		l.kind = kind
		l.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	default:
		err = errors.New("List of " + l.kind.String() + " kind.")
		return nil, err
	}

	l.count = 0
	l.pool = []*bucket{}
	if bucketSize > 0 {
		l.bucketSize = bucketSize
	} else {
		l.bucketSize = 1024
	}

	l.currentPushBucket = -1 // -1 means there is no available bucket in pool
	return
}

// Len returns the number of elements of list l.
// The complexity is O(1).
func (l *List) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.count
}

// Pop an element from list
func (l *List) Pop() (e Elem, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.count < 1 {
		err = errors.New("Pop from empty list")
		return
	}

	switch l.kind {
	case RAND:
		fallthrough
	case FIFO:
		bucket := l.pool[0]
		if bucket.popPosition < 0 && bucket.pushPosition > 0 {
			bucket.popPosition = 0
		}
		if bucket.popPosition >= 0 && bucket.popPosition < bucket.pushPosition {
			// Pop one element
			e = bucket.elems.Get(bucket.popPosition)
			bucket.popPosition++
			l.count--
		}

		// If first bucket is empty then pop empty bucket from pool buckets list
		if bucket.popPosition >= bucket.pushPosition {
			if len(l.pool) > 1 {
				l.pool = l.pool[1:]
				l.currentPushBucket--
			} else {
				l.pool = l.pool[:0]
				l.currentPushBucket = -1
			}
		}
		break
	case LIFO:
		bucket := l.pool[l.currentPushBucket]
		bucket.popPosition = bucket.pushPosition - 1
		if bucket.popPosition < 0 {
			l.pool = l.pool[:l.currentPushBucket]
			l.currentPushBucket--
			bucket = l.pool[l.currentPushBucket]
			bucket.popPosition = bucket.pushPosition - 1
		}
		// Pop one element
		e = bucket.elems.Get(bucket.popPosition)

		bucket.pushPosition--
		bucket.popPosition = bucket.pushPosition - 1

		l.count--

		break
	}

	return
}

// Push an element to list
func (l *List) Push(e Elem) (n int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.pool) == 0 {
		l.pool = append(l.pool, newBucket(l.bucketSize, e))
		l.currentPushBucket = 0
	}

	bucket := l.pool[l.currentPushBucket]
	if bucket.pushPosition == bucket.size {
		// Extend list when currentPushBucket is full
		l.pool = append(l.pool, newBucket(l.bucketSize, e))
		l.currentPushBucket++
		bucket = l.pool[l.currentPushBucket]
	}

	// shuffle
	if l.kind == RAND && l.count > 1 {
		// Random position
		p := l.r.Intn(l.count)
		// Align
		if l.pool[0].popPosition > 0 {
			p = p + l.pool[0].popPosition
		}
		// Location
		i := p / l.bucketSize
		m := p % l.bucketSize
		// Swap
		bucket.elems.Set(bucket.pushPosition, l.pool[i].elems.Get(m))
		l.pool[i].elems.Set(m, e)
	} else {
		// Push
		bucket.elems.Set(bucket.pushPosition, e)
	}

	// move position
	bucket.pushPosition++

	l.count++
	n = 1

	return
}
