package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Meas struct {
	min   float64
	max   float64
	sum   float64
	count float64
}

type Chunk struct {
	seek  int64
	limit int64
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing file")
	}

	filename := os.Args[1]
	chunks := rectifyChunk(filename)
	fmt.Println(chunks, len(chunks), cap(chunks))
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

	cpu := runtime.NumCPU()
	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	chunks := make([]Chunk, 0, cpu)

	size := stat.Size()
	chunk := size / int64(cpu)

	var prevChunk int64 = 0
	var nextChunk int64 = 0
	for nextChunk < size {
		seek := nextChunk + chunk
		if seek < size {
			_, err = file.Seek(seek, 0)
			if err != nil {
				log.Fatal(err)
			}

			buf, err := bufio.NewReader(file).ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			seek += int64(len(buf)) + 1
		}

		prevChunk = nextChunk
		nextChunk = seek
		if nextChunk > size {
			nextChunk = size
		}
		chunks = append(chunks, Chunk{seek: prevChunk, limit: nextChunk - prevChunk})
	}
	return chunks
}

func multiProcess(filename string, chunks []Chunk) {
	var wg sync.WaitGroup
	wg.Add(len(chunks))

	result := make([]map[string]*Meas, len(chunks))
	for i, chunk := range chunks {
		go func() {
			defer wg.Done()
			file, _ := os.Open(filename)
			defer func() {
				_ = file.Close()
			}()
			_, _ = file.Seek(chunk.seek, 0)

			meas := make(map[string]*Meas)
			scan(io.LimitReader(file, chunk.limit), meas)
			result[i] = meas
		}()
	}
	wg.Wait()

	meas := make(map[string]*Meas)
	for _, m := range result {
		for k, v := range m {
			if nv, ok := meas[k]; ok {
				nv.min = min(nv.min, v.min)
				nv.max = max(nv.max, v.max)
				nv.sum += v.sum
				nv.count += v.count
			} else {
				meas[k] = v
			}
		}
	}
	calcAndOutput(meas)
}

func scan(r io.Reader, meas map[string]*Meas) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, ";")
		temp, _ := strconv.ParseFloat(args[1], 64)

		v, ok := meas[args[0]]
		if !ok {
			v = &Meas{}
			meas[args[0]] = v
		}
		v.min = min(temp, v.min)
		v.max = max(temp, v.max)
		v.sum += temp
		v.count += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func calcAndOutput(meas map[string]*Meas) {
	keys := slices.Sorted(maps.Keys(meas))
	fmt.Print("{")
	for i, key := range keys {
		if i > 0 {
			fmt.Print(", ")
		}
		v := meas[key]
		fmt.Printf("%s=%.1f/%s/%.1f", key, v.min, round(v.sum/v.count), v.max)
	}
	fmt.Print("}")
}

func round(x float64) string {
	v := fmt.Sprintf("%.1f", x)
	if v == "-0.0" {
		v = "0.0"
	}
	return v
}
