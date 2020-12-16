package buffer

import (
	"errors"
	"io"
	"sync"
)

var ErrNodeEmpty = errors.New("empty")
var ErrBlockEmpty = errors.New("block empty")
var ErrBlockFull = errors.New("full")
var ErrBufferInputStreamHasFinish = errors.New("input stream finished")

type Block struct {
	data       []byte
	totalWrite int
	totalRead  int
	next       *Block
	blockSize  int
	mux        sync.RWMutex
}
type Buffer struct {
	pool           *sync.Pool
	head           *Block
	tail           *Block
	options        *Options
	inputStreamEnd bool
}

func NewBuffer(options ...Option) *Buffer {

	buffer := &Buffer{}
	buffer.options = loadOptions(options...)
	buffer.fix()
	blockSize := buffer.options.blockSize
	node := &Block{data: make([]byte, blockSize), blockSize: blockSize}
	buffer.head = node
	buffer.tail = node
	buffer.pool.New = func() interface{} {
		return &Block{data: make([]byte, blockSize), blockSize: blockSize}
	}
	return buffer
}
func (buffer *Buffer) fix() {
	if buffer.options.blockSize == 0 {
		buffer.options.blockSize = 4096
	}
	if buffer.options.pool != nil {
		buffer.pool = buffer.options.pool
	} else {
		buffer.pool = &sync.Pool{}
	}
}
func (buffer *Buffer) InputStreamFinish() {
	buffer.inputStreamEnd = true
}
func (buffer *Buffer) ReOpenInputStream() {
	buffer.inputStreamEnd = false
}

func (buffer *Buffer) Read(data []byte) (readLen int, err error) {
	readLen = 0
	defer func() {
		if err == nil && readLen == 0 {
			err = buffer.options.emptyError
			return
		}
	}()
	for buffer.head != nil && readLen < len(data) {
		buffer.head.mux.Lock()
		rd, err := buffer.head.Read(data[readLen:])
		readLen += rd
		if err != nil {
			if buffer.inputStreamEnd && buffer.head == buffer.tail && readLen == 0 {
				buffer.head.mux.Unlock()
				return 0, io.EOF
			}
			if err == ErrBlockEmpty {
				if buffer.head.next == nil {
					buffer.head.mux.Unlock()
					return readLen, nil
				}
				_head := buffer.head
				buffer.head = buffer.head.next
				_head.close()
				_head.mux.Unlock()
				buffer.pool.Put(_head)
				continue
			}
			if err == ErrNodeEmpty {
				buffer.head.mux.Unlock()
				return readLen, nil
			}
			return 0, err
		}

		buffer.head.mux.Unlock()
	}
	return readLen, nil
}
func (buffer *Buffer) Write(data []byte) (int, error) {
	if buffer.inputStreamEnd {
		return 0, ErrBufferInputStreamHasFinish
	}
	writeLen := 0
	for buffer.tail != nil && writeLen < len(data) {
		buffer.tail.mux.Lock()
		wd, err := buffer.tail.Write(data[writeLen:])
		if err != nil {
			if err == ErrBlockFull {
				newNode := buffer.pool.Get().(*Block)
				_tail := buffer.tail
				buffer.tail.next = newNode
				buffer.tail = newNode
				_tail.mux.Unlock()
				continue
			}
			return 0, err
		}
		writeLen += wd
		buffer.tail.mux.Unlock()
	}
	return writeLen, nil
}

func (bNode *Block) Read(data []byte) (int, error) {
	if bNode.totalRead == bNode.blockSize {
		return 0, ErrBlockEmpty
	}
	reset := bNode.totalWrite - bNode.totalRead
	if reset == 0 {
		return 0, ErrNodeEmpty
	}
	needRead := len(data)
	if reset > needRead {
		copy(data, bNode.data[bNode.totalRead:bNode.totalRead+needRead])
		bNode.totalRead += needRead
		return needRead, nil
	} else {
		copy(data, bNode.data[bNode.totalRead:bNode.totalWrite])
		bNode.totalRead += reset
		return reset, nil
	}
}
func (bNode *Block) Write(data []byte) (int, error) {
	needWrite := len(data)
	alreadyWriten := bNode.totalWrite
	resetWrite := bNode.blockSize - alreadyWriten
	if resetWrite == 0 {
		return 0, ErrBlockFull
	}
	if resetWrite > needWrite {
		copy(bNode.data[alreadyWriten:], data)
		bNode.totalWrite += needWrite
		return needWrite, nil
	} else {
		copy(bNode.data[alreadyWriten:], data[:resetWrite])
		bNode.totalWrite += resetWrite
		return resetWrite, nil
	}

}
func (bNode *Block) close() {
	bNode.totalRead = 0
	bNode.totalWrite = 0
	bNode.next = nil
	return
}
