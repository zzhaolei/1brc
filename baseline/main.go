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
)

type Meas struct {
	min   float64
	max   float64
	sum   float64
	count float64
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing file")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Close file err: %v\n", err)
		}
	}()

	meas := make(map[string]*Meas)
	scan(file, meas)
	calcAndOutput(meas)
}

func scan(r io.Reader, mea map[string]*Meas) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, ";")
		temp, _ := strconv.ParseFloat(args[1], 64)

		v, ok := mea[args[0]]
		if !ok {
			v = &Meas{}
			mea[args[0]] = v
		}
		v.min = min(temp, v.min)
		v.max = max(temp, v.max)
		v.sum = v.sum + temp
		v.count = v.count + 1
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
