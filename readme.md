## stream buffer
### 描述
受字节跳动[pool link block](https://mp.weixin.qq.com/s/wSaJYg-HqnYY4SdLA2Zzaw) 设计思想启发，streamBuffer是一个字节流，可以一个读和一个写并发进行,更加快。
内部使用go原生pool的设计，更加节省内存，当读写流平衡的时候，几乎不会消耗额外的内存。
### benchmark
通过benchmark可以看到，几乎不会额外申请内存。
```shell script
goos: darwin
goarch: amd64
pkg: github.com/gosundy/streambuffer
BenchmarkWrite
BenchmarkWrite-8                                    	    9277	    110645 ns/op	  101954 B/op	      48 allocs/op
BenchmarkRead
BenchmarkRead-8                                     	    2691	    428084 ns/op	     388 B/op	       0 allocs/op
BenchmarkWrite_Read_Concurrent_Ratio_Equal
BenchmarkWrite_Read_Concurrent_Ratio_Equal-8        	   44872	     25746 ns/op	      33 B/op	       3 allocs/op
BenchmarkWrite_Read_Concurrent_Ratio_Write_Fast
BenchmarkWrite_Read_Concurrent_Ratio_Write_Fast-8   	   61831	     20225 ns/op	      32 B/op	       3 allocs/op
BenchmarkWrite_Read_Concurrent_Ratio_Read_Fast
BenchmarkWrite_Read_Concurrent_Ratio_Read_Fast-8    	   50928	     23007 ns/op	      32 B/op	       3 allocs/op
PASS
```
### Usage
```go
func Demo(){
    blockSize := 4096
	buffer := NewBuffer(WithBlockSize(blockSize))
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
			if writeLn == totalWrite {
				break
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			rd, _ := buffer.Read(readData)
			readLn += rd
			if readLn == totalWrite {
				break
			}
		}
	}()
	wg.Wait()
}
```