package openkeyp2p

import (
	"container/list"
	"context"
	"sync"
)

// ThreadSafeContextBuffer ist ein threadsicherer Puffer mit Context-Unterstützung
type ThreadSafeContextBuffer struct {
	mu     sync.Mutex
	cond   *sync.Cond
	buffer *list.List
	ctx    context.Context
	cancel context.CancelFunc
	closed bool
}

// NewThreadSafeContextBuffer erstellt einen neuen Puffer mit Context
func NewThreadSafeContextBuffer(ctx context.Context) *ThreadSafeContextBuffer {
	childCtx, cancel := context.WithCancel(ctx)
	tscb := &ThreadSafeContextBuffer{
		buffer: list.New(),
		ctx:    childCtx,
		cancel: cancel,
	}
	tscb.cond = sync.NewCond(&tscb.mu)
	return tscb
}

// Put fügt Daten am Ende des Puffers hinzu
func (tscb *ThreadSafeContextBuffer) Put(data interface{}) error {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()

	select {
	case <-tscb.ctx.Done():
		return tscb.ctx.Err()
	default:
	}

	if tscb.closed {
		return context.Canceled
	}

	tscb.buffer.PushBack(data)
	tscb.cond.Signal()
	return nil
}

// Get wartet auf Daten und entfernt sie vom Anfang des Puffers
func (tscb *ThreadSafeContextBuffer) Get() (interface{}, error) {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()

	for tscb.buffer.Len() == 0 && !tscb.closed {
		select {
		case <-tscb.ctx.Done():
			return nil, tscb.ctx.Err()
		default:
			tscb.cond.Wait()
		}
	}

	if tscb.closed && tscb.buffer.Len() == 0 {
		return nil, context.Canceled
	}

	front := tscb.buffer.Front()
	data := front.Value
	tscb.buffer.Remove(front)
	return data, nil
}

// Prepend fügt Daten am Anfang des Puffers hinzu
func (tscb *ThreadSafeContextBuffer) Prepend(data interface{}) error {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()

	select {
	case <-tscb.ctx.Done():
		return tscb.ctx.Err()
	default:
	}

	if tscb.closed {
		return context.Canceled
	}

	tscb.buffer.PushFront(data)
	tscb.cond.Signal()
	return nil
}

// Close schließt den Puffer und weckt alle wartenden Leser
func (tscb *ThreadSafeContextBuffer) Close() {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()

	tscb.closed = true
	tscb.cancel()
	tscb.cond.Broadcast()
}

// Len gibt die aktuelle Anzahl der Elemente im Puffer zurück
func (tscb *ThreadSafeContextBuffer) Len() int {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()
	return tscb.buffer.Len()
}

// IsClosed prüft, ob der Puffer geschlossen wurde
func (tscb *ThreadSafeContextBuffer) IsClosed() bool {
	tscb.mu.Lock()
	defer tscb.mu.Unlock()
	return tscb.closed
}
