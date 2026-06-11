# Design & Architecture: Vol-Aware Indicators Engine

## 1. Ambiguity Decisions
* **IV Rank (Flat Window):** When $\text{Max} == \text{Min}$ (flat volatility), the indicator returns `0.0`. This prevents division-by-zero panics, `NaN`, or `+Inf` values, safely defaulting to a neutral baseline.
* **IV Percentile (Ties & Denominator):** * **Ties:** Handled using **strict inequality (`<`)**. Equal historical observations are excluded from the ranking count.
    * **Denominator:** Includes the current bar, keeping the total sample size strictly equal to the configured lookback `N`. 
    * *Note: This specific configuration was validated by matching the target `want 75` on Bar 3 in the golden verification test suite.*
* **Put-Call Ratio (PCR Market-Structure):**
    * **Zero-OI Regimes:** If a snapshot records zero open interest across all options strikes, `TotalCallOI` settles at zero. To prevent runtime division-by-zero failures, the engine intercepts this case and cleanly defaults to `0.0`.
    * **Inclusion Boundary:** The algorithm aggregates **all strikes** dynamically packaged inside the incoming `Snapshot.Chain` vector, capturing complete systemic flows without fragile price-band filtering.

## 2. Streaming Design & Complexity
The indicators leverage memory-bounded state management models:
* **Put-Call Ratio (PCR):** Processes snapshots point-in-time with $\mathcal{O}(1)$ space and $\mathcal{O}(K)$ time (where $K$ is the number of option strikes on the bar). It warms up instantly (`ready: true` on bar 1).
* **IV Percentile (IVP):** Implements a circular ring buffer with $\mathcal{O}(N)$ computation time.
* **IV Rank (Monotonic Deque Optimization):** Tracks historical window extrema inline during updates using dual monotonic double-ended queues.
    * **Time Complexity:** * `Update()`: Amortized $\mathcal{O}(1)$. Elements are added to the back, while stale values or sub-optimal extremes are popped off in an amortized fashion.
        * `Value()`: True $\mathcal{O}(1)$. Requires no loops; values are accessed directly via index pointer lookup (`r.minDeque[0].Value`).
    * **Space Complexity:** Strictly bounded at $\mathcal{O}(N)$ allocations.

## 3. Lookahead Prevention
The implementation is mathematically decoupled from lookahead bugs via three layers of isolation:
1. **Passive Execution:** Indicators are purely reactive; they only update when the `Evaluator` explicitly feeds them a historical close snapshot.
2. **Harness Timing:** The backtest harness drives updates sequentially on the *close* of bar $N$, while the execution engine only evaluates signals afterward.
3. **State Isolation:** The `resolveOperand` method reads the currently cached scalar state from `.Value()`. It has no reference or structural access to future array indices, making it impossible to peek at the open price of bar $N+1$.

## 4. Testing & Verification
* **Completed Tests:** Unit, golden file, and end-to-end integration tests pass cleanly, confirming proper type parsing from JSON types and correct trade triggers (`3 trades, net_pnl=-409.84`).
* **Next-Step Testing (Property-Based):** I would add an invariant validation test confirming that the streaming ring-buffer state exactly equals a batch-recomputed calculation (`incremental results == batch recompute`) over a shifting randomized sequence of 10,000 bars.

## 5. Production Tradeoffs
* **Time Budget vs. Optimization:** While IVP remains $\mathcal{O}(N)$ for window scans, IV Rank was optimized to an amortized $\mathcal{O}(1)$ approach using monotonic deques. This shows the ideal sweet spot for high-throughput live tickers without over-engineering components that don't scale out.
* **Production Enhancements:** In a live trading system, I would swap `float64` for fixed-point arithmetic (`shopspring/decimal`) to eliminate floating-point rounding accumulation errors across multi-year data horizons.