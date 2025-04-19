// Package snowflake provides a very simple Twitter snowflake generator.
package snowflake

import (
	"math/bits"
	"strconv"
	"sync"
	"time"
)

var (
	// DefaultEpoch is the default snowflake start timestamp epoch of 2020-01-01 00:00:00 UTC in milliseconds.
	DefaultEpoch int64 = 1577836800000
)

// Node is a struct holds the basic information needed for a snowflake generator
type Node struct {
	epoch int64 // start timestamp unix epoch in milliseconds
	node  int64

	mu    sync.Mutex
	start time.Time
	time  int64
	step  int64

	nodeMask  int64
	stepMask  int64
	timeShift int
	nodeShift int
}

// NewNode returns a new snowflake node that can be used to generate snowflake ID.
// The default epoch: 2000-01-01 00:00:00 UTC
// The default ID format:
// +----------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp | 10 Bit NodeID | 12 Bit Sequence ID |
// +----------------------------------------------------------------------+
// @see https://en.wikipedia.org/wiki/Snowflake_ID
func NewNode(node int64) *Node {
	return CustomNode(DefaultEpoch, node, 10, 12)
}

// CustomNode create a custom snowflake node.
func CustomNode(epoch, node int64, nodeBits, stepBits int) *Node {
	if nodeBits+stepBits > 22 {
		panic("snowflake: node+step bits must less than or equal 22")
	}

	nodeMax := ^(int64(-1) << nodeBits)
	if node < 0 || node > nodeMax {
		panic("snowflake: node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	n := Node{}
	n.epoch = epoch
	n.node = node
	n.nodeMask = nodeMax
	n.stepMask = ^(int64(-1) << stepBits)
	n.timeShift = nodeBits + stepBits
	n.nodeShift = stepBits

	now := time.Now()

	// add time.Duration to now to make sure we use the monotonic clock if available
	n.start = now.Add(time.Unix(epoch/1000, (epoch%1000)*1000000).Sub(now))

	return &n
}

// Epoch return the epoch
func (n *Node) Epoch() int64 {
	return n.epoch
}

// Node return the node
func (n *Node) Node() int64 {
	return n.node
}

// NodeBits return node bits
func (n *Node) NodeBits() int {
	return bits.Len64(uint64(n.nodeMask))
}

// StepBits return step bits
func (n *Node) StepBits() int {
	return n.nodeShift
}

// NextID creates and returns a unique snowflake ID
// To help guarantee uniqueness
// - Make sure your system is keeping accurate system time
// - Make sure you never have multiple nodes running with the same node ID
func (n *Node) NextID() ID {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.start).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.start).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	return ID{n, (now << n.timeShift) | (n.node << n.nodeShift) | n.step}
}

// LastID returns last createed snowflake ID
func (n *Node) LastID() ID {
	return ID{n, (n.time << n.timeShift) | (n.node << n.nodeShift) | n.step}
}

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID struct {
	sn *Node
	id int64
}

// Int64 returns an int64 of the snowflake ID
func (i ID) Int64() int64 {
	return i.id
}

// Time returns the snowflake ID time
func (i ID) Time() time.Time {
	return time.UnixMilli(i.UnixMilli())
}

// UnixMilli returns an int64 unix timestamp in milliseconds of the snowflake ID time
func (i ID) UnixMilli() int64 {
	return (i.id >> i.sn.timeShift) + i.sn.epoch
}

// Node returns an int64 of the snowflake ID node number
func (i ID) Node() int64 {
	return (i.id >> i.sn.nodeShift) & i.sn.nodeMask
}

// Step returns an int64 of the snowflake step (or sequence) number
func (i ID) Step() int64 {
	return i.id & i.sn.stepMask
}
