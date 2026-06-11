# Primary Task Prompts

1. content of README.md
2. content of indicators/ivrank.go
3. content of indicators/ivp.go
4. content of engine/evaluator.go
5. Based on the README goals and the 3 Go files provided, please:
    * Verify if the mathematical logic in `ivrank.go` and `ivp.go` perfectly aligns with the project specs.
    * Ensure `evaluator.go` handles data boundaries, empty sets, and errors correctly.

---

# Secondary Task Prompts

1. Implement the following:
    * **PCR (Put-Call Ratio):** `chain.csv` includes an `oi` column. Add a `PCR` indicator = total put OI / total call OI per bar, expose it in the DSL, and write a strategy combining it with `IV_RANK`. Note your market-structure assumptions (zero-OI strikes, which strikes to include).
    * **O(1) rolling extremes:** maintain IV Rank's sliding-window min/max in amortized O(1) with a monotonic deque instead of scanning the window each bar. Show the win with a benchmark.

    Let me know what you would require for this? 

2. content of market/market.go

3. First few rows data/chain.csv with headers

    ```csv
    time,strike,type,mark_iv,oi
    2025-01-01T00:00:00Z,58000.0000,CALL,0.6450,238.3193
    2025-01-01T00:00:00Z,58000.0000,PUT,0.6450,450.4316
    2025-01-01T00:00:00Z,59000.0000,CALL,0.6250,238.4446
    2025-01-01T00:00:00Z,59000.0000,PUT,0.6250,381.6017
    2025-01-01T00:00:00Z,60000.0000,CALL,0.6050,264.6376
    ```

4. content of indicators/indicator.go

5. content of indicators/ivrank.go

6. content of indicators/ivrank_test.go