package main

import (
	"fmt"
	"math"
	"math/rand"
)

type StorageEngine struct {
}

func main() {
	n := 100000
	hll := NewWithErrorRate(0.02)
	set := make(map[string]bool)

	fmt.Println("Inserted Values:", n)
	fmt.Println("Target Error Rate:", 0.02)
	r := rand.New(rand.NewSource(42))
	for i := 0; i < n; i++ {
		num := fmt.Sprintf("%d", r.Intn(50000))

		if ok := set[num]; !ok {
			set[num] = true
		}

		hll.Add(num)
	}
	exactCount := len(set)

	sum := 0
	for _, v := range hll.buckets {
		sum += v
	}
	error := math.Abs(float64(hll.Count())-float64(exactCount)) / float64(exactCount) * 100.0
	fmt.Println(hll.Count(), " vs ", exactCount, error, "%")
}
