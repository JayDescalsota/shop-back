package idgen

import (
	"sync"
	"time"
)

const (
	epoch        int64 = 1716768000000
	nodeBits           = 10
	stepBits           = 12
	nodeMax           int64 = -1 ^ (-1 << nodeBits)
	stepMax           int64 = -1 ^ (-1 << stepBits)
	timeShift         uint8 = nodeBits + stepBits
	nodeShift         uint8 = stepBits
)

type Snowflake struct {
	mu        sync.Mutex
	nodeID    int64
	step      int64
	lastTime  int64
}

func NewSnowflake(nodeID int64) (*Snowflake, error) {
	if nodeID < 0 || nodeID > nodeMax {
		return nil, ErrNodeOutOfRange
	}
	return &Snowflake{nodeID: nodeID}, nil
}

var ErrNodeOutOfRange = NewError("node id out of range")

var ErrRuntime = NewError("runtime error")

type IDError struct {
	msg string
}

func NewError(msg string) *IDError {
	return &IDError{msg: msg}
}

func (e *IDError) Error() string {
	return e.msg
}

func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli() - epoch

	if now == s.lastTime {
		s.step = (s.step + 1) & stepMax
		if s.step == 0 {
			for now <= s.lastTime {
				now = time.Now().UnixMilli() - epoch
			}
		}
	} else {
		s.step = 0
	}

	s.lastTime = now
	return (now << timeShift) | (s.nodeID << nodeShift) | s.step
}

func GenerateNodeID() int64 {
	return time.Now().UnixMilli() % nodeMax
}
