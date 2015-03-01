// Package hashtable is designed to cache large amount of data and avoid long GC pasue problem cause by using golang map too much
//
package hashtable

import "sync"

type Hashtable struct {
	count           int         // Total number of elements. Read only
	size            int         // Hash index size. Set when initial, read only other time
	lock            *sync.Mutex // Global lock
	index           []elemChain // Hash index array,
	traversePos     int         // current course position of traverse elements by step
	traverseKeysPos int         // current course position of traverse elements key by step
}

// Chain to save elements has same key hash value
type elemChain struct {
	elems Elems         // Elements
	lock  *sync.RWMutex // Chain lock
}
