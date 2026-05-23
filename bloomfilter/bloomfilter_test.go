package bloomfilter_test

import (
	"fmt"
	"testing"

	"beanstato.dev/logstrata/bloomfilter"
)

func TestBloomFilter_AddAndCheck(t *testing.T) {
	bf := bloomfilter.New(100000, 0.02)

	items := []string{"red", "green", "blue"}
	for _, item := range items {
		bf.Add(item)
	}

	// Should exist
	for _, item := range items {
		if !bf.Check(item) {
			t.Errorf("expected %q to be found, but Check returned false", item)
		}
	}
}

func TestBloomFilter_CheckAbsentItems(t *testing.T) {
	bf := bloomfilter.New(100000, 0.02)

	bf.Add("red")
	bf.Add("green")
	bf.Add("blue")

	absent := []string{"re", "gree", "bl", "yellow", "purple"}
	falsePositives := 0
	for _, item := range absent {
		if bf.Check(item) {
			falsePositives++
			t.Logf("false positive: %q reported as present", item)
		}
	}

	// With n=100000 and errorRate=0.02, false positives on a tiny set should be rare
	if falsePositives > 1 {
		t.Errorf("too many false positives: %d out of %d", falsePositives, len(absent))
	}
}

func TestBloomFilter_FalsePositiveRate(t *testing.T) {
	n := 100000
	errorRate := 0.02
	bf := bloomfilter.New(n, errorRate)

	// Fill with n items
	for i := 0; i < n; i++ {
		bf.Add(fmt.Sprintf("item-%d", i))
	}

	// Check items that were never added
	total := 10000
	falsePositives := 0
	for i := n; i < n+total; i++ {
		if bf.Check(fmt.Sprintf("item-%d", i)) {
			falsePositives++
		}
	}

	actualRate := float64(falsePositives) / float64(total)

	if actualRate > errorRate*2 {
		t.Errorf("false positive rate too high: got %.4f, want <= %.4f", actualRate, errorRate*2)
	}
	t.Logf("false positive rate: %.4f (target: %.4f)", actualRate, errorRate)
}
