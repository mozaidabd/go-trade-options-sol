package indicators

import "takehome-vol-indicators/market"

// DequeEntry bundles an IV value and the unique sequence ID it occurred at,
// allowing us to purge expired window entries in O(1).
type DequeEntry struct {
	Index uint64
	Value float64
}

// IVRank tracks where the current ATM IV sits within its trailing high/low range
// over the last `lookback` bars using an amortized O(1) sliding monotonic deque system.
type IVRank struct {
	lookback int
	window   []float64 // Holds rolling ATM IV values for reference
	head     int       // Pointer for circular ring buffer
	count    int       // Tracks warm-up status up to lookback capacity
	sequence uint64    // Monotonically increasing sequence counter per Update

	// Monotonic deques for tracking extrema
	minDeque []DequeEntry // Strictly increasing values
	maxDeque []DequeEntry // Strictly decreasing values
}

func NewIVRank(lookback int) *IVRank {
	return &IVRank{
		lookback: lookback,
		window:   make([]float64, lookback),
		head:     0,
		count:    0,
		sequence: 0,
		minDeque: make([]DequeEntry, 0, lookback),
		maxDeque: make([]DequeEntry, 0, lookback),
	}
}

func (r *IVRank) Name() string { return "IV_RANK" }

func (r *IVRank) Update(s market.Snapshot) {
	if r.lookback <= 0 {
		return
	}

	currentIV := s.ATMIV
	currentSeq := r.sequence
	r.sequence++

	// 1. Maintain Circular Ring Buffer
	r.window[r.head] = currentIV
	r.head = (r.head + 1) % r.lookback

	if r.count < r.lookback {
		r.count++
	}

	// 2. Expire elements outside our rolling historical sequence threshold
	if currentSeq >= uint64(r.lookback) {
		expirySeq := currentSeq - uint64(r.lookback)
		if len(r.minDeque) > 0 && r.minDeque[0].Index <= expirySeq {
			r.minDeque = r.minDeque[1:]
		}
		if len(r.maxDeque) > 0 && r.maxDeque[0].Index <= expirySeq {
			r.maxDeque = r.maxDeque[1:]
		}
	}

	// 3. Maintain Min Deque (Pop larger elements from back)
	for len(r.minDeque) > 0 && r.minDeque[len(r.minDeque)-1].Value >= currentIV {
		r.minDeque = r.minDeque[:len(r.minDeque)-1]
	}
	r.minDeque = append(r.minDeque, DequeEntry{Index: currentSeq, Value: currentIV})

	// 4. Maintain Max Deque (Pop smaller elements from back)
	for len(r.maxDeque) > 0 && r.maxDeque[len(r.maxDeque)-1].Value <= currentIV {
		r.maxDeque = r.maxDeque[:len(r.maxDeque)-1]
	}
	r.maxDeque = append(r.maxDeque, DequeEntry{Index: currentSeq, Value: currentIV})
}

func (r *IVRank) Value() (float64, bool) {
	// WARM-UP: Fully blocked until the window tracks enough sequence bars
	if r.count < r.lookback {
		return 0, false
	}

	// O(1) Fetch: The extrema are always positioned at the front of our deques
	minIV := r.minDeque[0].Value
	maxIV := r.maxDeque[0].Value

	// FLAT WINDOW: Safe guard against division-by-zero
	if maxIV == minIV {
		return 0.0, true
	}

	// Fetch current IV from ring buffer history
	currentIdx := (r.head - 1 + r.lookback) % r.lookback
	currentIV := r.window[currentIdx]

	ivRank := 100.0 * (currentIV - minIV) / (maxIV - minIV)
	return ivRank, true
}
