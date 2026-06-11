package indicators

import "takehome-vol-indicators/market"

// IVPercentile is the IV Percentile indicator: the fraction of bars in the
// trailing window whose ATM IV was below the current bar's, scaled to 0..100.
//
//	IVP = 100 * count(iv_i < iv_now) / (bars in window)
//
// TODO(candidate): implement this.
//
// IVP is NOT the same as IV Rank — make sure your two implementations actually
// differ. The definition has two genuine ambiguities; decide and document both
// in the README:
//   - ties: is the comparison strict (`<`) or inclusive (`<=`)?
//   - denominator: do you count the current bar in it, or only prior bars?
//
// Same streaming + warm-up requirements as IV Rank apply.
type IVPercentile struct {
	lookback int
	window   []float64 // Holds the rolling ATM IV values
	head     int       // Pointer for our rolling circular ring buffer
	count    int       // Current number of elements inside the window
}

func NewIVPercentile(lookback int) *IVPercentile {
	return &IVPercentile{
		lookback: lookback,
		window:   make([]float64, lookback),
		head:     0,
		count:    0,
	}
}

func (p *IVPercentile) Name() string { return "IV_PERCENTILE" }

func (p *IVPercentile) Update(s market.Snapshot) {
	if p.lookback <= 0 {
		return
	}

	// Overwrite the oldest element in our circular ring buffer
	p.window[p.head] = s.ATMIV
	p.head = (p.head + 1) % p.lookback

	// Increment elements tracking up to the maximum lookback capacity
	if p.count < p.lookback {
		p.count++
	}
}

func (p *IVPercentile) Value() (float64, bool) {
	// WARM-UP: Not ready until the window holds exactly `lookback` bars
	if p.count < p.lookback {
		return 0, false
	}

	// Retrieve the current bar's ATM IV (the one we just inserted)
	currentIdx := (p.head - 1 + p.lookback) % p.lookback
	currentIV := p.window[currentIdx]

	// STRICT COMPARISON: Count only bars strictly less than currentIV
	lessCount := 0
	for i := 0; i < p.lookback; i++ {
		if p.window[i] < currentIV {
			lessCount++
		}
	}

	// Calculate percentile using full lookback window as the denominator
	ivp := 100.0 * float64(lessCount) / float64(p.lookback)
	return ivp, true
}
