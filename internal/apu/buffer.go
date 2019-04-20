package apu

import "sync"

func newBuffer(capacity int) *buffer {
	return &buffer{
		bytes:    make([]byte, capacity),
		capacity: capacity,
	}
}

type buffer struct {
	mu                   sync.RWMutex
	bytes                []byte
	start, end, capacity int
}

func (b *buffer) size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := b.end - b.start
	if result < 0 {
		result += b.capacity
	}
	return result
}

func (b *buffer) push(vals []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if the current push is expected to overflow the capacity
	overflow := b.size()+len(vals) >= b.capacity

	n1 := copy(b.bytes[b.end:], vals)
	b.end += n1
	b.end %= b.capacity
	if n1 < len(vals) {
		// Wrapped around the end of the buffer
		n2 := copy(b.bytes, vals[n1:])
		b.end += n2
		b.end %= b.capacity
	}

	if overflow {
		// Buffer was overrun - make start right after end so that we use the full buffer
		b.start = (b.end + 1) % b.capacity
	}
}

func (b *buffer) Read(out []byte) (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.end > b.start {
		n := copy(out, b.bytes[b.start:b.end])
		b.start += n
		b.start %= b.capacity
		return n, nil
	}

	// Buffer read may cross circular bounds - copy is a little more complex
	n1 := copy(out, b.bytes[b.start:])
	b.start += n1
	b.start %= b.capacity
	if n1 == len(out) {
		return n1, nil
	}
	n2 := copy(out[n1:], b.bytes[:b.end])
	b.start += n2
	b.start %= b.capacity
	return n1 + n2, nil
}
