# Prompt Rationale Breakdown

## Primary Tasks Prompts Rationale

* **Prompt 1 (`README.md`)**: Feeds the AI the project's source of truth and architectural definitions. This establishes the constraints and formulas required to evaluate subsequent code steps accurately.
* **Prompts 2 & 3 (`indicators/ivrank.go` & `ivp.go`)**: Isolates the core mathematical files as standalone inputs. This prompts the AI to focus heavily on checking edge-case formulas, potential precision loss, and data boundaries without drowning in full-repo context.
* **Prompt 4 (`engine/evaluator.go`)**: Directs the AI to analyze the runtime execution layer. This allows the prompt sequence to transition from static math validation to operational stability testing (handling empty sets, array bounds, and errors).
* **Prompt 5 (Verification Objective)**: Synthesizes the previous four inputs. It explicitly forces the AI to cross-reference the math files (Prompts 2 & 3) against the specs (Prompt 1), and audit runtime safety (Prompt 4).

---

## Secondary Tasks Prompts Rationale

* **Prompt 1 (Task Specifications)**: Declares the engineering requirements and constraints early. It alerts the AI to look for data-aggregation paradigms (for PCR) and algorithmic performance targets (for the monotonic deque).
* **Prompt 2 (`market/market.go`)**: Supplies the AI with the exact structure of the data pipeline. This lets the AI map out how option chain rows are bundled into sequential time bars before writing any new aggregation logic.
* **Prompt 3 (`data/chain.csv`)**: Provides a raw schema sample. This explicitly guides the AI on string matching (`"CALL"`/`"PUT"`) and data parsing required for the PCR implementation.
* **Prompt 4 (`indicators/indicator.go`)**: Exposes the extension interface to the AI. This ensures that when the AI generates the PCR code, it integrates natively into the existing DSL structure.
* **Prompt 5 (`indicators/ivrank.go`)**: Feeds the target implementation for optimization into the context window. This gives the AI the exact logic it is tasked to rewrite into an amortized $O(1)$ algorithm.
* **Prompt 6 (`indicators/ivrank_test.go`)**: Feeds the testing matrix to the AI. This prompts the AI to keep test parity intact (regression testing) and provides the skeleton needed to generate an accurate Go performance benchmark.