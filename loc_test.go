package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestLocate(t *testing.T) {

}

func BenchmarkLocate(b *testing.B) {
	b.StopTimer()

	p := NewLocation()

	rand.Seed(time.Now().UnixNano())
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = p.Locate("1381580****")
	}
}
