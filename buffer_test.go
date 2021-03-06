package buffer

import (
	"sync"
	"testing"
)

const (
	TestBlockSize = 4096
)

func TestBPoolNode_Write(t *testing.T) {

	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	for i := 1; i <= 100; i++ {
		testData := make([]byte, i*buffer.options.blockSize-1)
		writeLn, _ := buffer.Write(testData)
		if writeLn != len(testData) {
			t.Fatalf("expect %d, actual:%d", len(testData), writeLn)
		}
	}
}
func TestBPoolNode_Read(t *testing.T) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	for i := 1; i <= 100; i++ {
		testData := make([]byte, i*buffer.options.blockSize-1)
		writeLn, _ := buffer.Write(testData)
		if writeLn != len(testData) {
			t.Fatalf("expect %d, actual:%d", len(testData), writeLn)
		}
		buffer.InputStreamFinish()
		totalRead := 0
		for {
			readData := make([]byte, 100)
			readln, err := buffer.Read(readData)
			totalRead += readln
			if err != nil {
				break
			}
		}
		if totalRead != writeLn {
			t.Fatalf("expect read:%d, acutal:%d", writeLn, totalRead)
		}
		buffer.ReOpenInputStream()
	}
}
func TestBPoolNode_IsEmpty(t *testing.T) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	for i := 1; i <= 100; i++ {
		testData := make([]byte, i*buffer.options.blockSize-1)
		writeLn, _ := buffer.Write(testData)
		if writeLn != len(testData) {
			t.Fatalf("expect %d, actual:%d", len(testData), writeLn)
		}
		buffer.InputStreamFinish()
		totalRead := 0
		for {
			readData := make([]byte, 100)
			readln, err := buffer.Read(readData)
			totalRead += readln
			if err != nil {
				break
			}
		}
		if totalRead != writeLn {
			t.Fatalf("expect read:%d, acutal:%d", writeLn, totalRead)
		}
		buffer.ReOpenInputStream()
	}
	if !buffer.IsEmpty() {
		t.Fatalf("expect is empty, acutal not empty")
	}
}
func TestBPoolBuffer_Write_Read_Concurrent_Ratio_Equal(t *testing.T) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	writeData := make([]byte, 100)
	readData := make([]byte, 100)

	for i := 0; i < 100; i++ {
		totalWrite := 10000000
		writeLn := 0
		readLn := 0
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			for {
				wd, _ := buffer.Write(writeData)
				writeLn += wd
				//t.Logf("write:%d",writeLn)
				if writeLn == totalWrite {
					//t.Log("write complete")
					break
				}
				//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
			}
		}()
		go func() {
			defer wg.Done()
			for {
				rd, _ := buffer.Read(readData)
				readLn += rd
				//	t.Logf("readln:%d",readLn)
				if readLn == totalWrite {
					//t.Log("read complete")
					break
				}
				//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
			}
		}()
		wg.Wait()
	}
}
func TestBPoolBuffer_Write_Read_Concurrent_Ratio_Write_Fast(t *testing.T) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	totalWrite := 100000
	writeLn := 0
	readLn := 0
	writeData := make([]byte, 1000)
	readData := make([]byte, 100)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			wd, _ := buffer.Write(writeData)
			writeLn += wd
			//	t.Logf("write:%d",writeLn)
			if writeLn == totalWrite {
				//	t.Log("write complete")
				break

			}
			//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			rd, _ := buffer.Read(readData)
			readLn += rd
			//t.Logf("readln:%d", readLn)
			if readLn == totalWrite {
				//t.Log("read complete")
				break
			}
			//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
		}
	}()
	wg.Wait()
}

func TestBPoolBuffer_Write_Read_Concurrent_Ratio_Read_Fast(t *testing.T) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	totalWrite := 100000
	writeLn := 0
	readLn := 0
	writeData := make([]byte, 100)
	readData := make([]byte, 1000)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			wd, _ := buffer.Write(writeData)
			writeLn += wd
			//	t.Logf("write:%d",writeLn)
			if writeLn == totalWrite {
				//t.Log("write complete")
				break

			}
			//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			rd, _ := buffer.Read(readData)
			readLn += rd
			//t.Logf("readln:%d",readLn)
			if readLn == totalWrite {
				//t.Log("read complete")
				break
			}
			//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
		}
	}()
	wg.Wait()
}
