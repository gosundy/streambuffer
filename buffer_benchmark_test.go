package buffer

import (
	"sync"
	"testing"
)

func BenchmarkWrite(b *testing.B) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		totalWrite := 100000
		writeLn := 0
		writeData := make([]byte, 100)
		for {
			wd, _ := buffer.Write(writeData)
			writeLn += wd
			if writeLn == totalWrite {
				break
			}
		}
	}

}
func BenchmarkRead(b *testing.B) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		totalWrite := 1000000
		writeLn := 0
		writeData := make([]byte, 100)
		for {
			wd, _ := buffer.Write(writeData)
			writeLn += wd
			if writeLn == totalWrite {
				break
			}
		}
		totalRead := 0
		for {
			readData := make([]byte, 1000)
			readln, _ := buffer.Read(readData)
			totalRead += readln
			//b.Log(totalRead)
			if writeLn == totalRead {
				break
			}

		}
	}
}
func BenchmarkWrite_Read_Concurrent_Ratio_Equal(b *testing.B) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	writeData := make([]byte, 100)
	readData := make([]byte, 100)
	totalWrite := 10000
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}

		writeLn := 0
		readLn := 0

		wg.Add(2)
		go func() {
			defer wg.Done()
			for {
				wd, _ := buffer.Write(writeData)
				writeLn += wd
				//t.Logf("write:%d",writeLn)
				if writeLn == totalWrite {
					//b.Log("write complete")
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
				if readLn == totalWrite {
					//b.Log("read complete")
					break
				}

				//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
			}
		}()
		wg.Wait()
	}

}
func BenchmarkWrite_Read_Concurrent_Ratio_Write_Fast(b *testing.B) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	writeData := make([]byte, 1000)
	readData := make([]byte, 100)
	totalWrite := 10000
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}

		writeLn := 0
		readLn := 0

		wg.Add(2)
		go func() {
			defer wg.Done()
			for {
				wd, _ := buffer.Write(writeData)
				writeLn += wd
				//t.Logf("write:%d",writeLn)
				if writeLn == totalWrite {
					//b.Log("write complete")
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
				if readLn == totalWrite {
					//b.Log("read complete")
					break
				}

				//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
			}
		}()
		wg.Wait()
	}

}
func BenchmarkWrite_Read_Concurrent_Ratio_Read_Fast(b *testing.B) {
	buffer := NewBuffer(WithBlockSize(TestBlockSize))
	writeData := make([]byte, 100)
	readData := make([]byte, 1000)
	totalWrite := 10000
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}

		writeLn := 0
		readLn := 0

		wg.Add(2)
		go func() {
			defer wg.Done()
			for {
				wd, _ := buffer.Write(writeData)
				writeLn += wd
				//t.Logf("write:%d",writeLn)
				if writeLn == totalWrite {
					//b.Log("write complete")
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
				if readLn == totalWrite {
					//b.Log("read complete")
					break
				}

				//time.Sleep(time.Duration(rand.Int31n(10)) * time.Millisecond)
			}
		}()
		wg.Wait()
	}

}
