package indicators

import "takehome-vol-indicators/market"

// IVRank is the IV Rank indicator: where the current ATM IV sits within its
// trailing high/low range over the last `lookback` bars, scaled to 0..100:
//
//	IV Rank = 100 * (iv_now - min) / (max - min)   over the trailing window
//
// TODO(candidate): implement this.
//
// Requirements:
//   - STREAMING: Update must do bounded work per bar (maintain a rolling window),
//     not rescan all history ever seen.
//   - WARM-UP: not ready until the window holds `lookback` bars; Value's bool is
//     false until then.
//   - FLAT WINDOW: decide what to return when max == min (no range) and document
//     your choice in the README. Don't return NaN/Inf or panic.
type IVRank struct {
	lookback int
	window   []float64 // Holds the rolling ATM IV values
	head     int       // Pointer for our rolling circular ring buffer
	count    int       // Current number of elements inside the window
}

func NewIVRank(lookback int) *IVRank {
	return &IVRank{
		lookback: lookback,
		window:   make([]float64, lookback),
		head:     0,
		count:    0,
	}
}

func (r *IVRank) Name() string { return "IV_RANK" }

func (r *IVRank) Update(s market.Snapshot) {
	if r.lookback <= 0 {
		return
	}

	// Overwrite the oldest element in our circular ring buffer
	r.window[r.head] = s.ATMIV
	r.head = (r.head + 1) % r.lookback

	// Increment elements tracking up to the maximum lookback capacity
	if r.count < r.lookback {
		r.count++
	}
}

func (r *IVRank) Value() (float64, bool) {
	// WARM-UP: Not ready until the window holds exactly `lookback` bars
	if r.count < r.lookback {
		return 0, false
	}

	// Find min and max within our fixed rolling window bounds
	minIV := r.window[0]
	maxIV := r.window[0]

	for i := 1; i < r.lookback; i++ {
		val := r.window[i]
		if val < minIV {
			minIV = val
		}
		if val > maxIV {
			maxIV = val
		}
	}

	// FLAT WINDOW: When max == min, there is no range (volatility is perfectly flat).
	// We return 0.0 to safely avoid NaN/Inf division-by-zero errors.
	if maxIV == minIV {
		return 0.0, true
	}

	// The current IV is the one we *just* inserted, which sits right behind the updated head pointer.
	currentIdx := (r.head - 1 + r.lookback) % r.lookback
	currentIV := r.window[currentIdx]

	ivRank := 100.0 * (currentIV - minIV) / (maxIV - minIV)
	return ivRank, true
}
