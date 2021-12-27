package ccu

import "github.com/pandafw/pango/col"

// UnboundedChan unbounded channel
type UnboundedChan struct {
	In  chan<- T        // channel for write
	Out <-chan T        // channel for read
	rb  *col.RingBuffer // buffer
}

// NewUnboundedChan create a unbounded channel
func NewUnboundedChan(initCapacity int) UnboundedChan {
	in := make(chan T, initCapacity)
	out := make(chan T, initCapacity)
	rb := col.NewRingBuffer(initCapacity)

	ubc := UnboundedChan{
		In:  in,
		Out: out,
		rb:  rb,
	}

	// start a goroutine, read from IN and write to OUT or BUFFER
	go func() {
		// close 'out' channel, after go-routine finish
		defer close(out)

	loop:
		for {
			val, ok := <-in
			if !ok { // 'in' is closed
				break loop
			}

			// try to write val to out channel
			select {
			case out <- val:
				// write ok
				continue
			default:
			}

			// out channel is full, write to buffer
			rb.Push(val)

			// drain the ring buffer
			for !rb.IsEmpty() {
				select {
				case val, ok := <-in: // read from 'in' to prevent 'in' blocking
					if !ok { // 'in' is closed
						break loop
					}
					rb.Push(val)
				case out <- rb.MustPeek(): // remove and get head data, write to the out channel
					rb.Poll()
				}
			}
		}

		// drain ring buffer, after 'in' in closed
		for !rb.IsEmpty() {
			out <- rb.MustPoll()
		}
	}()

	return ubc
}

// Len returns len of Out plus len of buffer.
func (ubc UnboundedChan) Len() int {
	return ubc.rb.Len() + len(ubc.Out)
}

// BufLen returns len of the buffer.
func (ubc UnboundedChan) BufLen() int {
	return ubc.rb.Len()
}

// IsEmpty returns true if the container length == 0
func (ubc UnboundedChan) IsEmpty() bool {
	return ubc.Len() == 0
}

// Clear clears the container
func (ubc UnboundedChan) Clear() {
	for {
		if _, ok := ubc.Poll(); !ok {
			return
		}
	}
}

// Push adds items of vs to the tail of queue
func (ubc UnboundedChan) Push(vs ...T) {
	for _, v := range vs {
		ubc.In <- v
	}
}

// Poll Retrieves and removes the head of this queue, or returns (nil, false) if this queue is empty.
func (ubc UnboundedChan) Poll() (v T, ok bool) {
	select {
	case v, ok = <-ubc.Out:
	default:
	}
	return v, ok
}

// Close Close the In channel
func (ubc UnboundedChan) Close() {
	close(ubc.In)
}
