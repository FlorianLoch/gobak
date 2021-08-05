package internal

import (
	"os"
	"sync"
)

type WorksetItem struct {
	path string
	info os.FileInfo
}

type Workset struct {
	mu sync.Mutex
	buf chan *WorksetItem
	set map[*WorksetItem]struct{}
}

func NewWorkset(bufferSize int) *Workset {
	return &Workset{
		buf: make(chan *WorksetItem, bufferSize),
		set: make(map[*WorksetItem]struct{}),
	}
}

func (w *Workset) Put(item *WorksetItem) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, in := w.set[item]; in {
		return
	}

	w.buf <- item

	w.set[item] = struct{}{}
}

func (w *Workset) C() <-chan *WorksetItem {
	ch := make(chan *WorksetItem) // unbuffered channel

	go func() {
		for item := range w.buf {
			w.mu.Lock()
			delete(w.set, item)
			w.mu.Unlock()

			ch <- item
		}
	}()

	return ch
}

func (w *Workset) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()

	return len(w.set)
}
