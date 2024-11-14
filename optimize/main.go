package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"runtime/pprof"
	"slices"
	"sync"
	"time"
)

const (
	Goroutine  = 25
	BufferSize = 4096 * 4096
)

type MeasMap map[uint64]*Meas
type MeasMap2 map[string]*Meas

type Meas struct {
	city  string
	min   int64
	max   int64
	sum   int64
	count int64
}

type Chunk struct {
	seek  int64
	limit int64
}

func main() {
	start := time.Now()

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = f.Close()
	}()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	run()
	_, _ = fmt.Fprintln(os.Stderr, time.Since(start).Milliseconds())
}

func run() {
	if len(os.Args) != 2 {
		log.Fatal("missing file")
	}

	filename := os.Args[1]
	chunks := rectifyChunk(filename)
	multiProcess(filename, chunks)
}

func rectifyChunk(filename string) []Chunk {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	chunks := make([]Chunk, 0, Goroutine)

	size := stat.Size()
	chunk := size / int64(Goroutine)

	var (
		prev int64 = 0
		next int64 = 0
	)
	for next < size {
		next += chunk
		if next < size {
			_, err = file.Seek(next, 0)
			if err != nil {
				log.Fatal(err)
			}

			buf, err := bufio.NewReaderSize(file, 128).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			next += int64(len(buf))
		} else {
			next = size
		}

		chunks = append(chunks, Chunk{seek: prev, limit: next - prev})
		prev = next
	}
	return chunks
}

func multiProcess(filename string, chunks []Chunk) {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	result := make([]MeasMap, len(chunks))
	for i, chunk := range chunks {
		go func() {
			defer wg.Done()
			file, _ := os.Open(filename)
			defer func() {
				_ = file.Close()
			}()
			_, _ = file.Seek(chunk.seek, 0)

			meas := make(MeasMap)
			scan(io.LimitReader(file, chunk.limit), meas)
			result[i] = meas
		}()
	}
	wg.Wait()

	meas := make(MeasMap2)
	for _, m := range result {
		for _, v := range m {
			if nv, ok := meas[v.city]; ok {
				nv.min = min(nv.min, v.min)
				nv.max = max(nv.max, v.max)
				nv.sum += v.sum
				nv.count += v.count
			} else {
				meas[v.city] = v
			}
		}
	}
	calcAndOutput(meas)
}

func scan(r io.Reader, meas MeasMap) {
	buffer := make([]byte, BufferSize)
	remain := make([]byte, 0, BufferSize)
	for {
		n, _ := r.Read(buffer[:len(buffer)-len(remain)])
		if n == 0 {
			break
		}

		remain = append(remain, buffer[:n]...)

		var (
			cityByte []byte
			tempByte []byte
			next     bool
			newBuf   = remain
		)
		for {
			cityByte, tempByte, newBuf, next = parseLine(newBuf)
			if !next { // 没有下一行，退出循环重新读取
				copy(remain, newBuf)
				remain = remain[:len(newBuf)]
				break
			}

			key := hash(cityByte)
			temp := parseNumber(tempByte)

			v, ok := meas[key]
			if !ok {
				v = &Meas{city: string(cityByte)}
				meas[key] = v
			}
			v.min = min(temp, v.min)
			v.max = max(temp, v.max)
			v.sum += temp
			v.count += 1
		}
	}
}

func parseLine(buffer []byte) (city, temp []byte, buf []byte, next bool) {
	end := 0
	for i, b := range buffer {
		if b != '\n' {
			continue
		}

		next = true
		end = i
		break
	}

	if !next {
		buf = buffer
		return
	}

	idx := 0
	for i, b := range buffer[:end] {
		if b == ';' {
			idx = i
			break
		}
	}
	city = buffer[:idx]
	temp = buffer[idx+1 : end]
	buf = buffer[end+1:]
	return
}

func parseNumber(data []byte) int64 {
	var (
		result     int64
		isNegative bool
	)
	for _, b := range data {
		if b == '-' {
			isNegative = true
			continue
		}

		if b >= '0' && b <= '9' {
			result = result*10 + int64(b-'0')
		}
	}
	if isNegative {
		result = -result
	}
	return result
}

func hash(name []byte) uint64 {
	var h uint64 = 5381
	for _, b := range name {
		h = (h << 5) + h + uint64(b)
	}
	return h
}

func calcAndOutput(meas map[string]*Meas) {
	keys := slices.Sorted(maps.Keys(meas))
	fmt.Print("{")
	for i, key := range keys {
		if i > 0 {
			fmt.Print(", ")
		}
		v := meas[key]
		fmt.Printf("%s=%.1f/%s/%.1f", key, float64(v.min)/10, round((float64(v.sum)/10)/float64(v.count)), float64(v.max)/10)
	}
	fmt.Println("}")
}

func round(x float64) string {
	v := fmt.Sprintf("%.1f", x)
	if v == "-0.0" {
		v = "0.0"
	}
	return v
}
