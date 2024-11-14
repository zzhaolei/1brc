package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	Goroutine = 25
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
	start := time.Now()
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

			buf, err := bufio.NewReader(file).ReadBytes('\n')
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
		city := args[0]
		temp, _ := strconv.ParseFloat(args[1], 64)

		v, ok := meas[city]
		if !ok {
			v = &Meas{}
			meas[city] = v
		}
		v.min = min(temp, v.min)
		v.max = max(temp, v.max)
		v.sum += temp
		v.count += 1
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
	fmt.Println("}")
}

func round(x float64) string {
	v := fmt.Sprintf("%.1f", x)
	if v == "-0.0" {
		v = "0.0"
	}
	return v
}
