package buffer

import (
	"errors"
	"sync"
)

var BpNodeEmpty = errors.New("empty")
var BpBlockEmpty = errors.New("block empty")
var BpNodeFull = errors.New("full")

const (
	BlockSize = 4096
)

type BPoolNode struct {
	data    []byte
	total   int
	read    int
	isClose bool
	next    *BPoolNode
	mux     sync.RWMutex
}
type BPoolBuffer struct {
	pool *sync.Pool
	head *BPoolNode
	tail *BPoolNode
}

func NewBuffer() *BPoolBuffer {
	node := &BPoolNode{data: make([]byte, BlockSize), isClose: false}
	buffer := &BPoolBuffer{pool: &sync.Pool{}, head: node, tail: node}
	buffer.pool.New = func() interface{} {
		return &BPoolNode{data: make([]byte, BlockSize), isClose: false}
	}
	return buffer
}

func (buffer *BPoolBuffer) Read(data []byte) (int, error) {

	readln := 0
	for buffer.head != nil && readln < len(data) {
		buffer.head.mux.Lock()
		rd, err := buffer.head.Read(data[readln:])
		readln += rd
		if err != nil {
			if err == BpBlockEmpty {
				if buffer.head.next == nil {
					buffer.head.mux.Unlock()
					return readln, nil
				}
				_head := buffer.head
				buffer.head = buffer.head.next
				_head.Close()
				_head.mux.Unlock()
				buffer.pool.Put(_head)
				continue
			}
			if err == BpNodeEmpty {
				buffer.head.mux.Unlock()
				return readln, nil
			}
			return 0, err
		}

		buffer.head.mux.Unlock()
	}
	return readln, nil
}
func (buffer *BPoolBuffer) Write(data []byte) (int, error) {
	writenLn := 0
	for buffer.tail != nil && writenLn < len(data) {
		buffer.tail.mux.Lock()
		wd, err := buffer.tail.Write(data[writenLn:])
		if err != nil {
			if err == BpNodeFull {
				newNode := buffer.pool.Get().(*BPoolNode)
				_tail := buffer.tail
				buffer.tail.next = newNode
				buffer.tail = newNode
				_tail.mux.Unlock()
				continue
			}
			return 0, err
		}
		writenLn += wd
		buffer.tail.mux.Unlock()
	}
	return writenLn, nil
}
func (buffer *BPoolBuffer) Fetch() *BPoolNode {
	return buffer.pool.Get().(*BPoolNode)
}
func (bNode *BPoolNode) Read(data []byte) (int, error) {
	if bNode.read == BlockSize {
		return 0, BpBlockEmpty
	}
	reset := bNode.total - bNode.read
	if reset == 0 {
		return 0, BpNodeEmpty
	}
	needRead := len(data)
	if reset > needRead {
		copy(data, bNode.data[bNode.read:bNode.read+needRead])
		bNode.read += needRead
		return needRead, nil
	} else {
		copy(data, bNode.data[bNode.read:bNode.total])
		bNode.read += reset
		return reset, nil
	}
}
func (bNode *BPoolNode) Write(data []byte) (int, error) {
	needWrite := len(data)
	alreadyWriten := bNode.total
	resetWrite := BlockSize - alreadyWriten
	if resetWrite == 0 {
		return 0, BpNodeFull
	}
	if resetWrite > needWrite {
		copy(bNode.data[alreadyWriten:], data)
		bNode.total += needWrite
		return needWrite, nil
	} else {
		copy(bNode.data[alreadyWriten:], data[:resetWrite])
		bNode.total += resetWrite
		return resetWrite, nil
	}

}
func (bNode *BPoolNode) Close() error {
	bNode.read = 0
	bNode.total = 0
	bNode.next = nil
	return nil
}
