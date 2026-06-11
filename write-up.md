# Design & Architecture: Vol-Aware Indicators

## 1. Ambiguity Decisions
* **IV Rank (Flat Window):** When $\text{Max} == \text{Min}$ (flat volatility), the indicator returns `0.0`. This prevents division-by-zero panics, `NaN`, or `+Inf` values, safely defaulting to a neutral baseline.
* **IV Percentile (Ties & Denominator):** * **Ties:** Handled using **strict inequality (`<`)**. Equal historical observations are excluded from the ranking count.
    * **Denominator:** Includes the current bar, keeping the total sample size strictly equal to the configured lookback `N`. 
    * *Note: This specific configuration was validated by matching the target `want 75` on Bar 3 in the golden verification test suite.*

## 2. Streaming Design & Complexity
The indicators implement a stateful **Circular Ring Buffer** using a pre-allocated slice of size `N`. 
* **Time Complexity:** * `Update()`: $\mathcal{O}(1)$. It overwrites the oldest element using a wrapping pointer (`head = (head + 1) % N`).
    * `Value()`: $\mathcal{O}(N)$. A single scan across the fixed window extracts metrics. Given typical financial windows ($N \le 252$), this approach avoids the overhead and edge-case fragility of dual-deque systems.
* **Space Complexity:** $\mathcal{O}(N)$ bounded memory. The underlying array is allocated exactly once during initialization and never appends or expands.
* **Warm-Up:** Both indicators maintain an internal counter. They strictly return `(_, false)` until the internal state has processed exactly `N` bars, preventing invalid trade signals.

## 3. Lookahead Prevention
The implementation is mathematically decoupled from lookahead bugs via three layers of isolation:
1. **Passive Execution:** Indicators are purely reactive; they only update when the `Evaluator` explicitly feeds them a historical close snapshot.
2. **Harness Timing:** The backtest harness drives updates sequentially on the *close* of bar $N$, while the execution engine only evaluates signals afterward.
3. **State Isolation:** The `resolveOperand` method reads the currently cached scalar state from `.Value()`. It has no reference or structural access to future array indices, making it impossible to peek at the open price of bar $N+1$.

## 4. Testing & Verification
* **Completed Tests:** Unit, golden file, and end-to-end integration tests pass cleanly, confirming proper type parsing from JSON types and correct trade triggers (`3 trades, net_pnl=-409.84`).
* **Next-Step Testing (Property-Based):** I would add an invariant validation test confirming that the streaming ring-buffer state exactly equals a batch-recomputed calculation (`incremental results == batch recompute`) over a shifting randomized sequence of 10,000 bars.

## 5. Production Tradeoffs
* **Time Budget vs. Optimization:** For windows where $N > 10,000$, I would optimize the $\mathcal{O}(N)$ search in `Value()` to an $\mathcal{O}(1)$ search by tracking running extrema using monotonic queues. For standard option lookbacks, the simpler linear scan reduces code footprint and code-path risk.
* **Production Enhancements:** In a live trading system, I would swap `float64` for fixed-point arithmetic (`shopspring/decimal`) to eliminate floating-point rounding accumulation errors across multi-year data horizons.