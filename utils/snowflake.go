package utils

import (
	"fmt"
	"sync"
	"time"
)

var (
	// Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
	Epoch int64 = 1288834974657
	// NodeBits holds the number of bits to use for Node
	NodeBits uint8 = 10
	// StepBits holds the number of bits to use for Step
	StepBits uint8 = 12

	nodeMax   int64 = -1 ^ (-1 << NodeBits)
	nodeMask        = nodeMax << StepBits
	stepMask  int64 = -1 ^ (-1 << StepBits)
	timeShift       = NodeBits + StepBits
	nodeShift       = StepBits
)

type UniqNode struct {
	mu    sync.Mutex
	epoch time.Time
	time  int64
	node  int64
	step  int64

	nodeMax   int64
	nodeMask  int64
	stepMask  int64
	timeShift uint8
	nodeShift uint8
}

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewUniqNode(node int64) *UniqNode {
	nodeMax = -1 ^ (-1 << NodeBits)
	nodeMask = nodeMax << StepBits
	stepMask = -1 ^ (-1 << StepBits)
	timeShift = NodeBits + StepBits
	nodeShift = StepBits

	n := UniqNode{}
	n.node = node
	n.nodeMax = -1 ^ (-1 << NodeBits)
	n.nodeMask = n.nodeMax << StepBits
	n.stepMask = -1 ^ (-1 << StepBits)
	n.timeShift = NodeBits + StepBits
	n.nodeShift = StepBits

	if n.node < 0 || n.node > n.nodeMax {
		return nil
	}

	var curTime = time.Now()
	// add time.Duration to curTime to make sure we use the monotonic clock if available
	n.epoch = curTime.Add(time.Unix(Epoch/1000, (Epoch%1000)*1000000).Sub(curTime))

	return &n
}

func (n *UniqNode) GenInt() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	now := time.Since(n.epoch).Nanoseconds() / 1000000

	if now == n.time {
		n.step = (n.step + 1) & n.stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Nanoseconds() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now
	return (now)<<n.timeShift | (n.node << n.nodeShift) | (n.step)
}

func (n *UniqNode) GenStr() string {
	code := n.GenInt()
	return fmt.Sprintf("%d", code)
}

//=========================================================
var GUniqNode *UniqNode
