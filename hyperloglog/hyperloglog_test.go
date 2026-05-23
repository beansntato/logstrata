package hyperloglog_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"beanstato.dev/logstrata/hyperloglog"
)

func TestHyperLogLog_ErrorRate(t *testing.T) {
	n := 100000
	targetErrorRate := 0.02

	hll := hyperloglog.NewWithErrorRate(targetErrorRate)
	set := make(map[string]bool)

	r := rand.New(rand.NewSource(42))
	for i := 0; i < n; i++ {
		num := fmt.Sprintf("%d", r.Intn(50000))
		if !set[num] {
			set[num] = true
		}
		hll.Add(num)
	}

	exactCount := len(set)
	estimated := hll.Count()
	errorRate := math.Abs(float64(estimated)-float64(exactCount)) / float64(exactCount)

	t.Logf("exact=%d  estimated=%d  error=%.4f%%", exactCount, estimated, errorRate*100)

	if errorRate > targetErrorRate*2 {
		t.Errorf("error rate too high: got %.4f%%, want <= %.4f%%",
			errorRate*100, targetErrorRate*2*100)
	}
}
