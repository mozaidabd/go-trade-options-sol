package indicators

import (
	"takehome-vol-indicators/market"
)

// PCR implements the Put-Call Ratio streaming indicator.
// PCR = Total Put Open Interest / Total Call Open Interest per bar snapshot.
type PCR struct {
	currentValue float64
	isReady      bool
}

// NewPCR constructs a new Put-Call Ratio indicator.
func NewPCR() *PCR {
	return &PCR{
		currentValue: 0.0,
		isReady:      false,
	}
}

// Update processes a single bar snapshot, aggregating the open interest
// across all option chain contracts present on the bar.
func (p *PCR) Update(s market.Snapshot) {
	var totalCallOI float64
	var totalPutOI float64

	// Aggregate Open Interest grouped by type across all strikes
	for _, entry := range s.Chain {
		switch entry.Type {
		case market.Call:
			totalCallOI += entry.OI
		case market.Put:
			totalPutOI += entry.OI
		}
	}

	// Market-Structure Safeguard: If Call Open Interest is zero,
	// default the ratio to 0.0 to prevent division by zero or +Inf values.
	if totalCallOI == 0 {
		p.currentValue = 0.0
	} else {
		p.currentValue = totalPutOI / totalCallOI
	}

	// PCR is a structural point-in-time snapshot metric requiring zero lookback history
	p.isReady = true
}

// Value returns the current Put-Call Ratio.
func (p *PCR) Value() (float64, bool) {
	return p.currentValue, p.isReady
}

// Name identifies the indicator in the strategy DSL parser.
func (p *PCR) Name() string {
	return "PCR"
}
