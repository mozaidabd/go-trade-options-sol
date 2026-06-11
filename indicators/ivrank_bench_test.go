package indicators

import (
	"math/rand"
	"testing"
	"time"

	"takehome-vol-indicators/market"
)

// generateBenchmarkSnapshots builds a sequence of random snapshots to simulate live processing.
func generateBenchmarkSnapshots(numBars int) []market.Snapshot {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	out := make([]market.Snapshot, numBars)

	// Seed fixed random to keep benchmark runs deterministic
	rg := rand.New(rand.NewSource(42))

	for i := 0; i < numBars; i++ {
		out[i] = market.Snapshot{
			Time:  base.Add(time.Duration(i) * time.Hour),
			ATMIV: 10.0 + rg.Float64()*50.0, // IV between 10% and 60%
		}
	}
	return out
}

// BenchmarkIVRank_SmallWindow measures performance for a standard short-term trading window (N=14).
func BenchmarkIVRank_SmallWindow(b *testing.B) {
	snaps := generateBenchmarkSnapshots(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := NewIVRank(14)
		for _, s := range snaps {
			r.Update(s)
			_, _ = r.Value()
		}
	}
}

// BenchmarkIVRank_LargeWindow measures scaling performance for a long-term macro window (N=500).
// In an O(N) system, this will be significantly slower than SmallWindow.
// In an O(1) monotonic deque system, runtime will remain nearly flat.
func BenchmarkIVRank_LargeWindow(b *testing.B) {
	snaps := generateBenchmarkSnapshots(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := NewIVRank(500)
		for _, s := range snaps {
			r.Update(s)
			_, _ = r.Value()
		}
	}
}
